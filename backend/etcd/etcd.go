package etcd

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

const ServicePrefix = "/services/"

type Client struct {
	Client *clientv3.Client
	Logger *zap.Logger

	endpoints map[string]string
	mu        sync.RWMutex
}

type Config struct {
	Endpoints   []string `mapstructure:"endpoints"`
	Username    string   `mapstructure:"username"`
	Password    string   `mapstructure:"password"`
	DialTimeout int      `mapstructure:"dial_timeout"`
}

func New(cfg *Config, logger *zap.Logger) (*Client, error) {
	if len(cfg.Endpoints) == 0 {
		cfg.Endpoints = []string{"localhost:2379"}
	}

	if cfg.DialTimeout == 0 {
		cfg.DialTimeout = 5
	}

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   cfg.Endpoints,
		Username:    cfg.Username,
		Password:    cfg.Password,
		DialTimeout: time.Duration(cfg.DialTimeout) * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create etcd client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := cli.Status(ctx, cfg.Endpoints[0]); err != nil {
		logger.Warn("etcd server not available", zap.Error(err))
	} else {
		logger.Info("connected to etcd", zap.Strings("endpoints", cfg.Endpoints))
	}

	c := &Client{
		Client:    cli,
		Logger:    logger,
		endpoints: make(map[string]string),
	}

	go c.watchServices(context.Background())

	return c, nil
}

func (c *Client) ServiceKey(appName string) string {
	return ServicePrefix + appName
}

func (c *Client) GetEndpoint(appName string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	endpoint, ok := c.endpoints[appName]
	if !ok {
		return "", fmt.Errorf("service %s not found", appName)
	}
	return endpoint, nil
}

func (c *Client) RegisterService(ctx context.Context, appName, endpoint string) error {
	key := ServicePrefix + appName
	_, err := c.Client.Put(ctx, key, endpoint)
	if err == nil {
		c.mu.Lock()
		c.endpoints[appName] = endpoint
		c.mu.Unlock()
		c.Logger.Info("registered service", zap.String("app", appName), zap.String("endpoint", endpoint))
	}
	return err
}

func (c *Client) DeregisterService(ctx context.Context, appName string) error {
	key := ServicePrefix + appName
	_, err := c.Client.Delete(ctx, key)
	if err == nil {
		c.mu.Lock()
		delete(c.endpoints, appName)
		c.mu.Unlock()
		c.Logger.Info("deregistered service", zap.String("app", appName))
	}
	return err
}

func (c *Client) ListServices(ctx context.Context) (map[string]string, error) {
	resp, err := c.Client.Get(ctx, ServicePrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	services := make(map[string]string)
	for _, kv := range resp.Kvs {
		appName := strings.TrimPrefix(string(kv.Key), ServicePrefix)
		services[appName] = string(kv.Value)
	}

	c.mu.Lock()
	c.endpoints = services
	c.mu.Unlock()

	return services, nil
}

func (c *Client) watchServices(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			rch := c.Client.Watch(ctx, ServicePrefix, clientv3.WithPrefix())
			for {
				select {
				case <-ctx.Done():
					return
				case resp, ok := <-rch:
					if !ok {
						time.Sleep(time.Second)
						goto retry
					}
					if resp.Err() != nil {
						c.Logger.Error("etcd watch event error", zap.Error(resp.Err()))
						continue
					}

					for _, event := range resp.Events {
						switch event.Type {
						case clientv3.EventTypePut:
							appName := strings.TrimPrefix(string(event.Kv.Key), ServicePrefix)
							c.mu.Lock()
							c.endpoints[appName] = string(event.Kv.Value)
							c.mu.Unlock()
							c.Logger.Debug("service updated", zap.String("app", appName), zap.String("endpoint", string(event.Kv.Value)))

						case clientv3.EventTypeDelete:
							appName := strings.TrimPrefix(string(event.Kv.Key), ServicePrefix)
							c.mu.Lock()
							delete(c.endpoints, appName)
							c.mu.Unlock()
							c.Logger.Debug("service removed", zap.String("app", appName))
						}
					}
				}
			}
		retry:
		}
	}
}

func (c *Client) Put(ctx context.Context, key, value string) error {
	_, err := c.Client.Put(ctx, key, value)
	return err
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	resp, err := c.Client.Get(ctx, key)
	if err != nil {
		return "", err
	}
	if len(resp.Kvs) == 0 {
		return "", nil
	}
	return string(resp.Kvs[0].Value), nil
}

func (c *Client) Delete(ctx context.Context, key string) error {
	_, err := c.Client.Delete(ctx, key)
	return err
}

func (c *Client) PutWithTTL(ctx context.Context, key, value string, ttl time.Duration) error {
	leaseID, err := c.Client.Grant(ctx, int64(ttl.Seconds()))
	if err != nil {
		return err
	}
	_, err = c.Client.Put(ctx, key, value, clientv3.WithLease(leaseID.ID))
	return err
}

func (c *Client) Watch(ctx context.Context, key string) clientv3.Watcher {
	return clientv3.NewWatcher(c.Client)
}

func (c *Client) Close() error {
	return c.Client.Close()
}
