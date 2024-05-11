package resolver

import (
	"context"
	"emperror.dev/errors"
	pb "github.com/je4/miniresolver/v2/pkg/miniresolverproto"
	"github.com/je4/utils/v2/pkg/zLogger"
	"google.golang.org/grpc/resolver"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const RESOLVERSCHEMA = "miniresolver"

func NewMiniResolverResolverBuilder(miniResolverclient pb.MiniResolverClient, logger zLogger.ZLogger) resolver.Builder {
	return &miniResolverResolverBuilder{
		miniResolverclient: miniResolverclient,
		logger:             logger,
	}
}

type miniResolverResolverBuilder struct {
	miniResolverclient pb.MiniResolverClient
	logger             zLogger.ZLogger
}

func (mrrb *miniResolverResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &miniResolverResolver{
		target:             target,
		cc:                 cc,
		miniResolverclient: mrrb.miniResolverclient,
		logger:             mrrb.logger,
	}
	r.start()
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
}

func (r *miniResolverResolver) start() {
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
	r.cc.UpdateState(resolver.State{Addresses: addrs})
}
func (*miniResolverResolver) ResolveNow(o resolver.ResolveNowOptions) {}
func (*miniResolverResolver) Close()                                  {}
