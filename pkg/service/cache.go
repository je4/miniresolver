package service

import (
	"github.com/je4/utils/v2/pkg/zLogger"
	"golang.org/x/exp/maps"
	"sync"
	"time"
)

const minNextCallTimeout = 10 * time.Second

type serviceEntry struct {
	service   string
	addresses map[string]time.Time
	sort      []string
}

func (se *serviceEntry) nextCallTimeout() time.Duration {
	if len(se.sort) == 0 {
		return minNextCallTimeout
	}
	var m time.Duration = 1 * time.Hour
	for _, a := range se.addresses {
		x := -time.Until(a)
		if x < m {
			m = x
		}
	}
	m += 2 * time.Second
	if m < minNextCallTimeout {
		m = minNextCallTimeout
	}
	return m
}

func (se *serviceEntry) removeOld(timeout time.Duration) {
	newSort := make([]string, 0, len(se.sort))
	for _, addr := range se.sort {
		if svc, ok := se.addresses[addr]; ok {
			if time.Since(svc) < timeout {
				newSort = append(newSort, addr)
			} else {
				delete(se.addresses, addr)
			}
		}
	}
	se.sort = newSort
}

func (se *serviceEntry) headToTail() {
	if len(se.sort) <= 1 {
		return
	}
	a := se.sort[0]
	copy(se.sort, se.sort[1:])
	se.sort[len(se.sort)-1] = a
}

func (se *serviceEntry) addAddress(addr string) {
	if _, ok := se.addresses[addr]; !ok {
		se.sort = append(se.sort, "")
		copy(se.sort[1:], se.sort)
		se.sort[0] = addr
	}
	se.addresses[addr] = time.Now()
}

func (se *serviceEntry) removeAddress(addr string) {
	delete(se.addresses, addr)
	for i, a := range se.sort {
		if a == addr {
			se.sort = append(se.sort[:i], se.sort[i+1:]...)
			break
		}
	}
}

func (se *serviceEntry) getAddresses(timeout time.Duration) []string {
	se.removeOld(timeout)
	return se.sort
}

func (se *serviceEntry) getAddress(timeout time.Duration) (string, time.Duration) {
	se.removeOld(timeout)
	if len(se.sort) == 0 {
		return "", 10 * time.Second
	}
	a := se.sort[0]
	se.headToTail()
	nct := -time.Until(se.addresses[a])
	if nct < minNextCallTimeout {
		nct = minNextCallTimeout
	}
	return a, nct
}

func newCache(timeout time.Duration, logger zLogger.ZLogger) *cache {
	return &cache{
		Mutex:    sync.Mutex{},
		timeout:  timeout,
		services: make(map[string]*serviceEntry),
		logger:   logger,
	}
}

type cache struct {
	sync.Mutex
	timeout  time.Duration
	services map[string]*serviceEntry
	logger   zLogger.ZLogger
}

func (c *cache) addService(name, addr string, domains []string) {
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
		if !ok {
			svcs = &serviceEntry{
				service:   serviceName,
				addresses: make(map[string]time.Time),
				sort:      make([]string, 0, 1),
			}
			c.services[serviceName] = svcs
		}
		svcs.addAddress(addr)
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
