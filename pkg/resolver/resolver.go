package resolver

import (
	"context"
	"emperror.dev/errors"
	pb "github.com/je4/miniresolver/v2/pkg/miniresolverproto"
	"github.com/je4/utils/v2/pkg/zLogger"
	"google.golang.org/grpc/resolver"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"time"
)

const RESOLVERSCHEMA = "miniresolver"

func NewMiniResolverResolverBuilder(miniResolverclient pb.MiniResolverClient, checkTimeout time.Duration, notFoundTimeout time.Duration, logger zLogger.ZLogger) resolver.Builder {
	if time.Duration(checkTimeout).Seconds() == 0 {
		checkTimeout = 10 * time.Minute
	}
	if time.Duration(notFoundTimeout).Seconds() == 0 {
		notFoundTimeout = 10 * time.Second
	}
	return &miniResolverResolverBuilder{
		miniResolverclient: miniResolverclient,
		checkTimeout:       checkTimeout,
		notFoundTimeout:    notFoundTimeout,
		logger:             logger,
	}
}

type miniResolverResolverBuilder struct {
	miniResolverclient pb.MiniResolverClient
	logger             zLogger.ZLogger
	checkTimeout       time.Duration
	notFoundTimeout    time.Duration
}

func (mrrb *miniResolverResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &miniResolverResolver{
		target:             target,
		cc:                 cc,
		miniResolverclient: mrrb.miniResolverclient,
		logger:             mrrb.logger,
		done:               make(chan bool),
		checkTimeout:       mrrb.checkTimeout,
		notFoundTimeout:    mrrb.notFoundTimeout,
	}
	go func() {
		for {
			timeout := r.doIt()
			select {
			case <-r.done:
				return
			case <-time.After(timeout):
			}
		}
	}()
	return r, nil
}
func (*miniResolverResolverBuilder) Scheme() string { return RESOLVERSCHEMA }

// miniResolverResolver is a
// Resolver(https://godoc.org/google.golang.org/grpc/resolver#Resolver).
type miniResolverResolver struct {
	target             resolver.Target
	cc                 resolver.ClientConn
	miniResolverclient pb.MiniResolverClient
	logger             zLogger.ZLogger
	done               chan bool
	checkTimeout       time.Duration
	notFoundTimeout    time.Duration
}

func (r *miniResolverResolver) doIt() (timeout time.Duration) {
	timeout = r.notFoundTimeout
	addr := r.target.Endpoint()
	r.logger.Debug().Msgf("start resolver for %s", addr)
	resp, err := r.miniResolverclient.ResolveServices(context.Background(), &wrapperspb.StringValue{Value: addr})
	if err != nil {
		r.logger.Error().Err(err).Msgf("cannot resolve %s", addr)
		r.cc.ReportError(errors.Wrapf(err, "cannot resolve %s", addr))
		return
	}
	for _, a := range resp.Addr {
		r.logger.Debug().Msgf("resolved %s to '%s'", addr, a)
	}
	addrs := make([]resolver.Address, len(resp.Addr))
	for i, s := range resp.Addr {
		addrs[i] = resolver.Address{Addr: s}
	}
	if err := r.cc.UpdateState(resolver.State{Addresses: addrs}); err != nil {
		r.logger.Error().Err(err).Msgf("cannot update state for %s", addr)
		return
	}
	if len(resp.Addr) > 0 {
		timeout = r.checkTimeout
	}
	return
}
func (r *miniResolverResolver) ResolveNow(resolver.ResolveNowOptions) {
	//r.logger.Debug().Msgf("resolve now")
}
func (r *miniResolverResolver) Close() {
	r.logger.Debug().Msgf("close %s", r.target.Endpoint())
	r.done <- true
}
