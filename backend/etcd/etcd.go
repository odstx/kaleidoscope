package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

const ServicePrefix = "/services/"

type ServiceInstance struct {
	AppName      string            `json:"app_name"`
	Version      string            `json:"version"`
	InstanceID   string            `json:"instance_id"`
	Endpoint     string            `json:"endpoint"`
	Status       string            `json:"status"`
	Metadata     map[string]string `json:"metadata,omitempty"`
	RegisteredAt time.Time         `json:"registered_at"`
}

type ServiceInfo struct {
	AppName   string            `json:"app_name"`
	Version   string            `json:"version"`
	Instances []ServiceInstance `json:"instances"`
}

type Client struct {
	Client   *clientv3.Client
	Logger   *zap.Logger
	services map[string]map[string]map[string]ServiceInstance
	mu       sync.RWMutex
}

type Registry interface {
	RegisterService(ctx context.Context, appName, version, instanceID, endpoint string, metadata map[string]string) error
	DeregisterService(ctx context.Context, appName, version, instanceID string) error
	GetService(ctx context.Context, appName, version string) ([]ServiceInstance, error)
	GetServiceInstance(ctx context.Context, appName, version, instanceID string) (*ServiceInstance, error)
	ListServices(ctx context.Context) (map[string]ServiceInfo, error)
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
		Client:   cli,
		Logger:   logger,
		services: make(map[string]map[string]map[string]ServiceInstance),
	}

	go c.watchServices(context.Background())

	return c, nil
}

func (c *Client) serviceKey(appName, version, instanceID string) string {
	return fmt.Sprintf("%s%s/%s/%s", ServicePrefix, appName, version, instanceID)
}

func (c *Client) versionKey(appName, version string) string {
	return fmt.Sprintf("%s%s/%s", ServicePrefix, appName, version)
}

func (c *Client) RegisterService(ctx context.Context, appName, version, instanceID, endpoint string, metadata map[string]string) error {
	if version == "" {
		version = "v1"
	}
	if instanceID == "" {
		instanceID = fmt.Sprintf("%d", time.Now().UnixNano())
	}

	instance := ServiceInstance{
		AppName:      appName,
		Version:      version,
		InstanceID:   instanceID,
		Endpoint:     endpoint,
		Status:       "healthy",
		Metadata:     metadata,
		RegisteredAt: time.Now(),
	}

	data, err := json.Marshal(instance)
	if err != nil {
		return fmt.Errorf("failed to marshal service instance: %w", err)
	}

	key := c.serviceKey(appName, version, instanceID)
	_, err = c.Client.Put(ctx, key, string(data))
	if err != nil {
		return fmt.Errorf("failed to register service: %w", err)
	}

	c.mu.Lock()
	if c.services == nil {
		c.services = make(map[string]map[string]map[string]ServiceInstance)
	}
	if c.services[appName] == nil {
		c.services[appName] = make(map[string]map[string]ServiceInstance)
	}
	if c.services[appName][version] == nil {
		c.services[appName][version] = make(map[string]ServiceInstance)
	}
	c.services[appName][version][instanceID] = instance
	c.mu.Unlock()

	c.Logger.Info("registered service",
		zap.String("app", appName),
		zap.String("version", version),
		zap.String("instance_id", instanceID),
		zap.String("endpoint", endpoint))
	return nil
}

func (c *Client) DeregisterService(ctx context.Context, appName, version, instanceID string) error {
	if version == "" {
		version = "v1"
	}
	if instanceID == "" {
		return fmt.Errorf("instance_id is required")
	}

	key := c.serviceKey(appName, version, instanceID)
	_, err := c.Client.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to deregister service: %w", err)
	}

	c.mu.Lock()
	if c.services[appName] != nil && c.services[appName][version] != nil {
		delete(c.services[appName][version], instanceID)
	}
	c.mu.Unlock()

	c.Logger.Info("deregistered service",
		zap.String("app", appName),
		zap.String("version", version),
		zap.String("instance_id", instanceID))
	return nil
}

func (c *Client) GetService(ctx context.Context, appName, version string) ([]ServiceInstance, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.services == nil {
		return nil, nil
	}

	if appName == "" {
		var result []ServiceInstance
		for _, versions := range c.services {
			for _, instances := range versions {
				for _, instance := range instances {
					result = append(result, instance)
				}
			}
		}
		return result, nil
	}

	if version == "" {
		var result []ServiceInstance
		if versions, ok := c.services[appName]; ok {
			for _, instances := range versions {
				for _, instance := range instances {
					result = append(result, instance)
				}
			}
		}
		return result, nil
	}

	instances, ok := c.services[appName][version]
	if !ok {
		return nil, nil
	}

	result := make([]ServiceInstance, 0, len(instances))
	for _, instance := range instances {
		result = append(result, instance)
	}
	return result, nil
}

func (c *Client) GetServiceInstance(ctx context.Context, appName, version, instanceID string) (*ServiceInstance, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.services == nil {
		return nil, fmt.Errorf("service %s not found", appName)
	}

	if version == "" {
		version = "v1"
	}

	versions, ok := c.services[appName]
	if !ok {
		return nil, fmt.Errorf("service %s not found", appName)
	}

	instances, ok := versions[version]
	if !ok {
		return nil, fmt.Errorf("service %s version %s not found", appName, version)
	}

	instance, ok := instances[instanceID]
	if !ok {
		return nil, fmt.Errorf("instance %s not found", instanceID)
	}

	return &instance, nil
}

func (c *Client) GetHealthyInstance(ctx context.Context, appName, version string) (*ServiceInstance, error) {
	instances, err := c.GetService(ctx, appName, version)
	if err != nil {
		return nil, err
	}

	for i := range instances {
		if instances[i].Status == "healthy" {
			return &instances[i], nil
		}
	}

	return nil, fmt.Errorf("no healthy instance found for %s version %s", appName, version)
}

func (c *Client) ListServices(ctx context.Context) (map[string]ServiceInfo, error) {
	resp, err := c.Client.Get(ctx, ServicePrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	result := make(map[string]ServiceInfo)
	for _, kv := range resp.Kvs {
		var instance ServiceInstance
		if err := json.Unmarshal(kv.Value, &instance); err != nil {
			continue
		}

		key := string(kv.Key)
		parts := strings.Split(strings.TrimPrefix(key, ServicePrefix), "/")
		if len(parts) < 3 {
			continue
		}
		appName := parts[0]
		version := parts[1]

		info, ok := result[appName]
		if !ok {
			info = ServiceInfo{
				AppName: appName,
				Version: version,
			}
		}
		info.Instances = append(info.Instances, instance)
		result[appName] = info
	}

	c.mu.Lock()
	c.services = make(map[string]map[string]map[string]ServiceInstance)
	for appName, info := range result {
		if c.services[appName] == nil {
			c.services[appName] = make(map[string]map[string]ServiceInstance)
		}
		for _, instance := range info.Instances {
			version := instance.Version
			if c.services[appName][version] == nil {
				c.services[appName][version] = make(map[string]ServiceInstance)
			}
			c.services[appName][version][instance.InstanceID] = instance
		}
	}
	c.mu.Unlock()

	return result, nil
}

func (c *Client) ListAllServiceVersions(ctx context.Context) (map[string][]string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make(map[string][]string)
	for appName, versions := range c.services {
		for version := range versions {
			result[appName] = append(result[appName], version)
		}
	}
	return result, nil
}

func (c *Client) UpdateServiceStatus(ctx context.Context, appName, version, instanceID, status string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.services == nil {
		return fmt.Errorf("service %s not found", appName)
	}

	if version == "" {
		version = "v1"
	}

	versions, ok := c.services[appName]
	if !ok {
		return fmt.Errorf("service %s not found", appName)
	}

	instances, ok := versions[version]
	if !ok {
		return fmt.Errorf("service %s version %s not found", appName, version)
	}

	instance, ok := instances[instanceID]
	if !ok {
		return fmt.Errorf("instance %s not found", instanceID)
	}

	instance.Status = status
	c.services[appName][version][instanceID] = instance

	data, err := json.Marshal(instance)
	if err != nil {
		return fmt.Errorf("failed to marshal service instance: %w", err)
	}

	key := c.serviceKey(appName, version, instanceID)
	_, err = c.Client.Put(ctx, key, string(data))
	if err != nil {
		return fmt.Errorf("failed to update service status: %w", err)
	}

	c.Logger.Info("updated service status",
		zap.String("app", appName),
		zap.String("version", version),
		zap.String("instance_id", instanceID),
		zap.String("status", status))
	return nil
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
							var instance ServiceInstance
							if err := json.Unmarshal(event.Kv.Value, &instance); err != nil {
								continue
							}
							c.mu.Lock()
							if c.services == nil {
								c.services = make(map[string]map[string]map[string]ServiceInstance)
							}
							if c.services[instance.AppName] == nil {
								c.services[instance.AppName] = make(map[string]map[string]ServiceInstance)
							}
							if c.services[instance.AppName][instance.Version] == nil {
								c.services[instance.AppName][instance.Version] = make(map[string]ServiceInstance)
							}
							c.services[instance.AppName][instance.Version][instance.InstanceID] = instance
							c.mu.Unlock()
							c.Logger.Debug("service updated",
								zap.String("app", instance.AppName),
								zap.String("version", instance.Version),
								zap.String("instance_id", instance.InstanceID))

						case clientv3.EventTypeDelete:
							key := string(event.Kv.Key)
							parts := strings.Split(strings.TrimPrefix(key, ServicePrefix), "/")
							if len(parts) < 3 {
								continue
							}
							appName := parts[0]
							version := parts[1]
							instanceID := parts[2]
							c.mu.Lock()
							if c.services[appName] != nil && c.services[appName][version] != nil {
								delete(c.services[appName][version], instanceID)
							}
							c.mu.Unlock()
							c.Logger.Debug("service removed",
								zap.String("app", appName),
								zap.String("version", version),
								zap.String("instance_id", instanceID))
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
