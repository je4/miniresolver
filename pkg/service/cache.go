package service

import (
	"github.com/je4/utils/v2/pkg/zLogger"
	"golang.org/x/exp/maps"
	"sync"
	"time"
)

const minNextCallTimeout = 10 * time.Second

func newCache(timeout time.Duration, logger zLogger.ZLogger) *cache {
	c := &cache{
		Mutex:    sync.Mutex{},
		timeout:  timeout,
		services: make(map[string]*serviceEntry),
		logger:   logger,
		done:     make(chan bool),
	}
	c.Start()
	return c
}

type cache struct {
	sync.Mutex
	timeout  time.Duration
	services map[string]*serviceEntry
	logger   zLogger.ZLogger
	done     chan bool
}

func (c *cache) Close() {
	c.done <- true
	c.Lock()
	defer c.Unlock()
	for _, svcs := range c.services {
		svcs.Clear()
	}
}

func (c *cache) Start() {
	go func() {
		for {
			select {
			case <-time.After(time.Minute):
				c.removeUnavailable()
			case <-c.done:
				return
			}
		}
	}()
}

func (c *cache) addService(name, addr string, domains []string, single bool) {
	c.Lock()
	defer c.Unlock()
	if len(domains) == 0 {
		domains = []string{""}
	}
	for _, domain := range domains {
		serviceName := name
		if domain != "" {
			serviceName = domain + "." + name
		}
		svcs, ok := c.services[serviceName]
		if ok {
			if !svcs.refreshAddress(addr) {
				if single {
					svcs.Clear()
				}
				svcs.addAddress(addr)
			}
		} else {
			svcs = NewServiceEntry(serviceName, c.logger)
			svcs.addAddress(addr)
			c.services[serviceName] = svcs
		}

		c.logger.Debug().Msgf("service address added %s: %v", serviceName, addr)
	}
	for name, svcs := range c.services {
		if len(svcs.addresses) == 0 {
			delete(c.services, name)
		} else {
			c.logger.Debug().Msgf("service %s: %v in cache", name, maps.Keys(svcs.addresses))
		}
	}
}

func (c *cache) removeService(name, addr string, domains []string) {
	c.Lock()
	defer c.Unlock()
	if len(domains) == 0 {
		domains = []string{""}
	}
	for _, domain := range domains {
		if domain != "" {
			name = domain + "." + name
		}
		svcs, ok := c.services[name]
		if !ok {
			return
		}
		svcs.removeAddress(addr)
		if len(svcs.addresses) == 0 {
			delete(c.services, name)
		}
	}
}

func (c *cache) getServices(name string) ([]string, time.Duration) {
	c.Lock()
	defer c.Unlock()
	svcs, ok := c.services[name]
	if !ok {
		return []string{}, minNextCallTimeout
	}
	return svcs.getAddresses(c.timeout), svcs.nextCallTimeout()
}

func (c *cache) getService(name string) (string, time.Duration) {
	c.Lock()
	defer c.Unlock()
	svcs, ok := c.services[name]
	if !ok {
		return "", minNextCallTimeout
	}
	return svcs.getAddress(c.timeout)
}

func (c *cache) removeUnavailable() {
	for _, svcs := range c.services {
		addrs := svcs.getUnavailable()
		if len(addrs) > 0 {
			c.Lock()
			svcs.removeAddress(addrs...)
			c.Unlock()
		}
	}
}
