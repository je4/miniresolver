package service

import (
	"context"
	"fmt"
	pb "github.com/je4/miniresolver/v2/pkg/miniresolverproto"
	"github.com/je4/utils/v2/pkg/zLogger"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"net"
	"time"
)

func NewMiniResolver(bufferSize int, serviceExpiration time.Duration, logger zLogger.ZLogger) *MiniResolver {

	return &MiniResolver{
		logger:            logger,
		services:          newCache(serviceExpiration),
		serviceExpiration: serviceExpiration,
	}
}

type MiniResolver struct {
	pb.UnimplementedMiniResolverServer
	logger            zLogger.ZLogger
	services          *cache
	serviceExpiration time.Duration
}

func (d *MiniResolver) Ping(context.Context, *emptypb.Empty) (*pb.DefaultResponse, error) {
	return &pb.DefaultResponse{
		Status:  pb.ResultStatus_OK,
		Message: "pong",
		Data:    nil,
	}, nil
}

func (d *MiniResolver) AddService(ctx context.Context, data *pb.ServiceData) (*pb.DefaultResponse, error) {
	d.logger.Debug().Msgf("add service '%s' - '%s:%d'", data.GetService(), data.GetHost(), data.GetPort())

	var address = fmt.Sprintf("%s:%d", data.GetHost(), data.GetPort())
	if data.GetHost() == "" {
		p, ok := peer.FromContext(ctx)
		if !ok {
			return nil, fmt.Errorf("cannot get peer")
		}
		peerAddr := p.Addr.String()
		host, _, err := net.SplitHostPort(peerAddr)
		if err != nil {
			return nil, fmt.Errorf("cannot split host port of '%s': %v", peerAddr, err)
		}
		ip := net.ParseIP(host)
		if ip.To4() == nil {
			host = fmt.Sprintf("[%s]", host)
		}
		address = fmt.Sprintf("%s:%d", host, data.GetPort())
	}

	d.services.addService(data.Service, address)
	d.logger.Debug().Msgf("service '%s' - '%s' added", data.Service, address)
	return &pb.DefaultResponse{
		Status:  pb.ResultStatus_OK,
		Message: fmt.Sprintf("service '%s' - '%s' added", data.Service, address),
	}, nil
}

func (d *MiniResolver) RemoveService(ctx context.Context, data *pb.ServiceData) (*pb.DefaultResponse, error) {
	d.logger.Debug().Msgf("remove service '%s' - '%s:%d'", data.Service, data.GetHost(), data.GetPort())

	var address = fmt.Sprintf("%s:%d", data.GetHost(), data.GetPort())
	if data.GetHost() == "" {
		p, ok := peer.FromContext(ctx)
		if !ok {
			return nil, fmt.Errorf("cannot get peer")
		}
		peerAddr := p.Addr.String()
		host, _, err := net.SplitHostPort(peerAddr)
		if err != nil {
			return nil, fmt.Errorf("cannot split host port of '%s': %v", peerAddr, err)
		}
		ip := net.ParseIP(host)
		if ip.To4() == nil {
			host = fmt.Sprintf("[%s]", host)
		}
		address = fmt.Sprintf("%s:%d", host, data.GetPort())
	}
	d.services.removeService(data.Service, address)
	d.logger.Debug().Msgf("service '%s' - '%s' removed", data.Service, address)
	waitSeconds := int64((d.serviceExpiration.Seconds() * 2.0) / 3.0)
	return &pb.DefaultResponse{
		Status:      pb.ResultStatus_OK,
		Message:     fmt.Sprintf("service '%s' - '%s' removed", data.Service, address),
		WaitSeconds: &waitSeconds,
	}, nil
}

func (d *MiniResolver) ResolveServices(ctx context.Context, data *wrapperspb.StringValue) (*pb.ServiceResponse, error) {
	d.logger.Debug().Msgf("resolve services '%s'", data.Value)
	addrs := d.services.getServices(data.Value)
	return &pb.ServiceResponse{
		Addr: addrs,
	}, nil
}

func (d *MiniResolver) ResolveService(ctx context.Context, data *wrapperspb.StringValue) (*wrapperspb.StringValue, error) {
	d.logger.Debug().Msgf("resolve service '%s'", data.Value)
	addr := d.services.getService(data.Value)
	if addr == "" {
		return nil, fmt.Errorf("service '%s' not found", data.Value)
	}
	return &wrapperspb.StringValue{
		Value: addr,
	}, nil
}
