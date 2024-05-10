package service

import (
	"context"
	"github.com/bluele/gcache"
	pb "github.com/je4/miniresolver/v2/pkg/miniresolverproto"
	"github.com/je4/utils/v2/pkg/zLogger"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

func NewMiniResolver(bufferSize int, serviceExpiration time.Duration, logger zLogger.ZLogger) *MiniResolver {

	return &MiniResolver{
		logger:            logger,
		services:          gcache.New(bufferSize).Expiration(serviceExpiration).LRU().Build(),
		serviceExpiration: serviceExpiration,
	}
}

type MiniResolver struct {
	pb.UnimplementedMiniResolverServer
	logger            zLogger.ZLogger
	services          gcache.Cache
	serviceExpiration time.Duration
}

func (d *MiniResolver) Ping(context.Context, *emptypb.Empty) (*pb.DefaultResponse, error) {
	return &pb.DefaultResponse{
		Status:  pb.ResultStatus_OK,
		Message: "pong",
		Data:    nil,
	}, nil
}
