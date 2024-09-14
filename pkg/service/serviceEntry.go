package service

import (
	"context"
	"github.com/je4/utils/v2/pkg/zLogger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"time"
)

func NewServiceEntry(service string, logger zLogger.ZLogger) *serviceEntry {
	return &serviceEntry{
		service:   service,
		addresses: make(map[string]time.Time),
		client:    make(map[string]*grpc.ClientConn),
		sort:      make([]string, 0, 1),
		logger:    logger,
	}
}

type serviceEntry struct {
	service   string
	addresses map[string]time.Time
	client    map[string]*grpc.ClientConn
	sort      []string
	logger    zLogger.ZLogger
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
	olds := make([]string, 0, len(se.addresses))
	for _, addr := range se.sort {
		if svc, ok := se.addresses[addr]; ok {
			if time.Since(svc) >= timeout {
				olds = append(olds, addr)
			}
		}
	}
	se.removeAddress(olds...)
}

func (se *serviceEntry) getUnavailable() []string {
	unavailable := make([]string, 0, len(se.addresses))
	for addr := range se.addresses {
		cc, ok := se.client[addr]
		if !ok {
			se.logger.Debug().Msgf("ping %s::%s - client not found", se.service, addr)
			continue
		}
		if err := cc.Invoke(context.Background(), "ping", nil, nil); err != nil {
			stat := status.FromContextError(err)
			se.logger.Debug().Msgf("ping %s::%s - [%d]", se.service, addr, stat.Code())
			if stat.Code() == codes.Unavailable {
				se.logger.Debug().Msgf("%s::%s not available", se.service, addr)
				unavailable = append(unavailable, addr)
			}
		}

	}
	return unavailable
}

func (se *serviceEntry) headToTail() {
	if len(se.sort) <= 1 {
		return
	}
	a := se.sort[0]
	copy(se.sort, se.sort[1:])
	se.sort[len(se.sort)-1] = a
}
func (se *serviceEntry) refreshAddress(addr string) bool {
	if _, ok := se.addresses[addr]; ok {
		se.addresses[addr] = time.Now()
		return true
	}
	return false
}

func (se *serviceEntry) addAddress(addr string) {
	if _, ok := se.addresses[addr]; !ok {
		se.sort = append(se.sort, "")
		copy(se.sort[1:], se.sort)
		se.sort[0] = addr
	}
	se.addresses[addr] = time.Now()
	var err error
	se.client[addr], err = grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		delete(se.client, addr)
	}
	if err := se.client[addr].Invoke(context.Background(), "ping", nil, nil); err != nil {
		stat := status.FromContextError(err)
		if stat.Code() == codes.Unavailable {
			se.removeAddress(addr)
		}
	}
}

func (se *serviceEntry) removeAddress(addrs ...string) {
	for _, addr := range addrs {
		delete(se.addresses, addr)
		if c, ok := se.client[addr]; ok {
			c.Close()
			delete(se.client, addr)
		}
		for i, a := range se.sort {
			if a == addr {
				se.sort = append(se.sort[:i], se.sort[i+1:]...)
				break
			}
		}
	}
}

func (se *serviceEntry) Clear() {
	for _, c := range se.client {
		c.Close()
	}
	se.client = make(map[string]*grpc.ClientConn)
	se.addresses = make(map[string]time.Time)
	se.sort = make([]string, 0, 1)
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
