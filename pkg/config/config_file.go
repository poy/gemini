package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/poy/go-dependency-injection/pkg/injection"
	"github.com/poy/go-router/pkg/observability"
)

type KeyName string

type Config interface {
	Get(key string) string
	Set(key, value string)
	List() []string
	Store() error
}

type configFile struct {
	mu    sync.Mutex
	props map[string]string

	path string
}

func init() {
	injection.Register(func(ctx context.Context) Config {
		return loadConfig(ctx)
	})
}

func loadConfig(ctx context.Context) *configFile {
	log := injection.Resolve[observability.Logger](ctx)
	cfg := configFile{
		props: make(map[string]string),
	}
	defer cfg.setupProps(ctx)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("failed to read home directory: %v", err)
	}

	configDirPath := filepath.Join(homeDir, ".config", "gemini")

	if fi, err := os.Stat(configDirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(configDirPath, 0700); err != nil {
			log.Fatalf("failed to create directory: %v", err)
		}
	} else if err != nil {
		log.Fatalf("failed to stat config directory: %v", err)
	} else if !fi.IsDir() {
		log.Fatalf("expected %q to be a directory", configDirPath)
	}

	configFilePath := filepath.Join(configDirPath, "config.json")
	cfg.path = configFilePath

	configFile, err := os.Open(configFilePath)
	if os.IsNotExist(err) {
		return &cfg
	} else if err != nil {
		log.Fatalf("failed to open config file: %v", err)
	}
	defer configFile.Close()

	if err := json.NewDecoder(configFile).Decode(&cfg.props); err != nil {
		log.Fatalf("failed to decode config file: %v", err)
	}

	return &cfg
}

func (c *configFile) setupProps(ctx context.Context) {
	// Add any registered keys that don't have values and remove any non-registered ones.
	registeredKeys := injection.Resolve[injection.Group[KeyName]](ctx)
	registeredKeysM := make(map[string]bool)

	for _, k := range registeredKeys.Vals() {
		if _, ok := c.props[string(k)]; !ok {
			c.props[string(k)] = ""
		}
		registeredKeysM[string(k)] = true
	}
	var deleteKeys []string
	for k := range c.props {
		if registeredKeysM[k] {
			continue
		}
		deleteKeys = append(deleteKeys, k)
	}
	for _, k := range deleteKeys {
		delete(c.props, k)
	}
}

func (c *configFile) List() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	var keys []string
	for k := range c.props {
		keys = append(keys, k)
	}
	return keys
}

func (c *configFile) Get(key string) string {
	c.mu.Lock()
	defer c.mu.Unlock()
	val, ok := c.props[key]
	if !ok {
		panic(fmt.Sprintf("unregistered config key value: %s", key))
	}
	return val
}

func (c *configFile) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Ensure the key is registered
	_, ok := c.props[key]
	if !ok {
		panic(fmt.Sprintf("unregistered config key value: %s", key))
	}

	c.props[key] = value
}

func (c *configFile) Store() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	f, err := os.Create(c.path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(c.props); err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}
	return nil
}
