package middleware

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"kaleidoscope/config"
	"kaleidoscope/etcd"
	"kaleidoscope/models"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

func MicroserviceProxy(cfg *config.Config, db *gorm.DB, etcdClient *etcd.Client) gin.HandlerFunc {
	if !cfg.Microservice.Enabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	auth := CombinedAuth(cfg, db)
	tracer := otel.Tracer("kaleidoscope/proxy")

	return func(c *gin.Context) {
		path := c.Request.URL.Path

		if !strings.HasPrefix(path, "/app/") {
			c.Next()
			return
		}

		auth(c)
		if c.IsAborted() {
			return
		}

		userID, _ := c.Get("userID")
		uidStr := fmt.Sprintf("%v", userID)

		var username string
		var user models.User
		if err := db.Where("id = ?", uidStr).First(&user).Error; err == nil {
			username = user.Username
		}

		parts := strings.SplitN(strings.TrimPrefix(path, "/app/"), "/", 2)
		if len(parts) < 1 || parts[0] == "" {
			c.Next()
			return
		}

		appName := parts[0]
		version := ""
		targetPath := ""
		if len(parts) > 1 {
			targetPath = "/" + parts[1]
		}
		if len(parts) > 2 {
			version = parts[1]
			targetPath = "/" + parts[2]
		}

		versionHeader := c.GetHeader("X-Version")
		if versionHeader != "" {
			if versionHeader == "latest" {
				versions, err := etcdClient.ListAllServiceVersions(c.Request.Context())
				if err == nil && len(versions[appName]) > 0 {
					version = versions[appName][len(versions[appName])-1]
				}
			} else {
				version = versionHeader
			}
		}

		if version == "" {
			versions, err := etcdClient.ListAllServiceVersions(c.Request.Context())
			if err == nil && len(versions[appName]) > 0 {
				version = versions[appName][len(versions[appName])-1]
			}
		}

		if version == "" {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "service not found"})
			c.Abort()
			return
		}

		instance, err := etcdClient.GetHealthyInstance(c.Request.Context(), appName, version)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "service unavailable"})
			c.Abort()
			return
		}

		targetURL := fmt.Sprintf("http://%s%s", instance.Endpoint, targetPath)

		target, _ := url.Parse(targetURL)

		proxy := httputil.NewSingleHostReverseProxy(target)

		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {

			originalDirector(req)
			req.URL = target

			req.Header.Set("X-USER-UID", uidStr)
			req.Header.Set("X-SOURCE", "kaleidoscope")
			if username != "" {
				req.Header.Set("X-User-Name", username)
			}
			otel.GetTextMapPropagator().Inject(req.Context(), propagation.HeaderCarrier(req.Header))
		}

		ctx, span := tracer.Start(c.Request.Context(), "proxy to "+appName,
			trace.WithAttributes(
				attribute.String("app.name", appName),
				attribute.String("target.url", targetURL),
				attribute.String("user.id", uidStr),
			))
		defer span.End()

		c.Request = c.Request.WithContext(ctx)
		proxy.ServeHTTP(c.Writer, c.Request)
		c.Abort()
	}
}
