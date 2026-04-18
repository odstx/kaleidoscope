package etcd

import (
	"context"
	"net/http"
	"sync"
	"time"

	"kaleidoscope/metrics"

	"go.uber.org/zap"
)

type HealthChecker struct {
	client      *Client
	logger      *zap.Logger
	checkInterval time.Duration
	timeout       time.Duration
	stopCh        chan struct{}
	instanceHealth map[string]bool
	mu            sync.RWMutex
}

func NewHealthChecker(client *Client, logger *zap.Logger, checkInterval, timeout time.Duration) *HealthChecker {
	return &HealthChecker{
		client:        client,
		logger:        logger,
		checkInterval: checkInterval,
		timeout:       timeout,
		stopCh:        make(chan struct{}),
		instanceHealth: make(map[string]bool),
	}
}

func (hc *HealthChecker) Start() {
	hc.logger.Info("Starting health checker")
	go hc.run()
}

func (hc *HealthChecker) Stop() {
	close(hc.stopCh)
	hc.logger.Info("Health checker stopped")
}

func (hc *HealthChecker) run() {
	ticker := time.NewTicker(hc.checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-hc.stopCh:
			return
		case <-ticker.C:
			hc.checkAllInstances()
		}
	}
}

func (hc *HealthChecker) checkAllInstances() {
	ctx, cancel := context.WithTimeout(context.Background(), hc.timeout)
	defer cancel()

	services, err := hc.client.ListServices(ctx)
	if err != nil {
		hc.logger.Error("Failed to list services for health check", zap.Error(err))
		return
	}

	now := time.Now()
	for appName, info := range services {
		for _, instance := range info.Instances {
			healthy := hc.checkInstance(instance.Endpoint)
			instanceID := instance.InstanceID

			hc.mu.Lock()
			wasHealthy := hc.instanceHealth[appName+"/"+instanceID]
			hc.instanceHealth[appName+"/"+instanceID] = healthy
			hc.mu.Unlock()

			if healthy != wasHealthy || wasHealthy {
				metrics.RecordInstanceHealth(appName, instanceID, healthy)
			}

			if !healthy && wasHealthy {
				hc.logger.Warn("Instance became unhealthy",
					zap.String("app", appName),
					zap.String("instance_id", instanceID),
					zap.String("endpoint", instance.Endpoint))
				hc.updateInstanceStatus(ctx, appName, instance.Version, instanceID, "unhealthy")
			} else if healthy && !wasHealthy {
				hc.logger.Info("Instance became healthy",
					zap.String("app", appName),
					zap.String("instance_id", instanceID),
					zap.String("endpoint", instance.Endpoint))
				hc.updateInstanceStatus(ctx, appName, instance.Version, instanceID, "healthy")
			}

			hc.client.mu.Lock()
			if appServices, ok := hc.client.services[appName]; ok {
				if versionInstances, ok := appServices[instance.Version]; ok {
					if inst, ok := versionInstances[instanceID]; ok {
						inst.Status = "healthy"
						versionInstances[instanceID] = inst
					}
				}
			}
			hc.client.mu.Unlock()

			_ = now
		}
	}
}

func (hc *HealthChecker) checkInstance(endpoint string) bool {
	url := endpoint + "/health"
	if endpoint[len(endpoint)-1] == '/' {
		url = endpoint + "health"
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return false
	}

	req = req.WithContext(context.Background())

	client := &http.Client{Timeout: hc.timeout}
	resp, err := client.Do(req)
	if err != nil {
		hc.logger.Debug("Health check failed", zap.String("endpoint", endpoint), zap.Error(err))
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode >= 200 && resp.StatusCode < 300
}

func (hc *HealthChecker) updateInstanceStatus(ctx context.Context, appName, version, instanceID, status string) {
	err := hc.client.UpdateServiceStatus(ctx, appName, version, instanceID, status)
	if err != nil {
		hc.logger.Error("Failed to update instance status",
			zap.String("app", appName),
			zap.String("instance_id", instanceID),
			zap.String("status", status),
			zap.Error(err))
	}
}

func (hc *HealthChecker) GetInstanceHealth(appName, instanceID string) bool {
	hc.mu.RLock()
	defer hc.mu.RUnlock()
	return hc.instanceHealth[appName+"/"+instanceID]
}

func (hc *HealthChecker) IsHealthy(appName, version, instanceID string) bool {
	healthy := hc.checkInstanceFromCache(appName, instanceID)
	if !healthy {
		endpoint := hc.getInstanceEndpoint(appName, version, instanceID)
		if endpoint != "" {
			healthy = hc.checkInstance(endpoint)
			hc.mu.Lock()
			hc.instanceHealth[appName+"/"+instanceID] = healthy
			hc.mu.Unlock()
		}
	}
	return healthy
}

func (hc *HealthChecker) checkInstanceFromCache(appName, instanceID string) bool {
	hc.mu.RLock()
	defer hc.mu.RUnlock()
	health, ok := hc.instanceHealth[appName+"/"+instanceID]
	if !ok {
		return true
	}
	return health
}

func (hc *HealthChecker) getInstanceEndpoint(appName, version, instanceID string) string {
	hc.client.mu.RLock()
	defer hc.client.mu.RUnlock()

	if hc.client.services == nil {
		return ""
	}

	if appName == "" {
		return ""
	}

	versions, ok := hc.client.services[appName]
	if !ok {
		return ""
	}

	if version == "" {
		for _, v := range versions {
			if instance, ok := v[instanceID]; ok {
				return instance.Endpoint
			}
		}
		return ""
	}

	instances, ok := versions[version]
	if !ok {
		return ""
	}

	if instance, ok := instances[instanceID]; ok {
		return instance.Endpoint
	}

	return ""
}