package service

import (
	"sync"
	"time"
)

type serviceEntry struct {
	service   string
	addresses map[string]time.Time
	sort      []string
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

func (se *serviceEntry) getAddress(timeout time.Duration) string {
	if len(se.sort) == 0 {
		return ""
	}

	for _, addr := range se.sort {
		if svc, ok := se.addresses[addr]; ok {
			if time.Since(svc) < timeout {
				return addr
			}
		}
	}
	se.sort = make([]string, 0, 1)
	se.addresses = make(map[string]time.Time, 1)
	return ""
}

func newCache(timeout time.Duration) *cache {
	return &cache{
		Mutex:    sync.Mutex{},
		timeout:  timeout,
		services: make(map[string]*serviceEntry),
	}
}

type cache struct {
	sync.Mutex
	timeout  time.Duration
	services map[string]*serviceEntry
}

func (c *cache) addService(name, addr string) {
	c.Lock()
	defer c.Unlock()
	svcs, ok := c.services[name]
	if !ok {
		svcs = &serviceEntry{
			service:   name,
			addresses: make(map[string]time.Time),
			sort:      make([]string, 0, 1),
		}
		c.services[name] = svcs
	}
	svcs.addAddress(addr)
}

func (c *cache) removeService(name, addr string) {
	c.Lock()
	defer c.Unlock()
	svcs, ok := c.services[name]
	if !ok {
		return
	}
	svcs.removeAddress(addr)
	if len(svcs.addresses) == 0 {
		delete(c.services, name)
	}
}

func (c *cache) getServices(name string) []string {
	c.Lock()
	defer c.Unlock()
	svcs, ok := c.services[name]
	if !ok {
		return []string{}
	}
	return svcs.getAddresses(c.timeout)
}

func (c *cache) getService(name string) string {
	c.Lock()
	defer c.Unlock()
	svcs, ok := c.services[name]
	if !ok {
		return ""
	}
	return svcs.getAddress(c.timeout)
}
