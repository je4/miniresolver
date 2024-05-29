package resolver

import (
	"context"
	genericproto "github.com/je4/genericproto/v2/pkg/generic/proto"
	"github.com/je4/utils/v2/pkg/zLogger"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"reflect"
)

type GRPCPinger interface {
	Ping(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*genericproto.DefaultResponse, error)
}

func DoPing(pinger GRPCPinger, logger zLogger.ZLogger) {
	pingerType := reflect.TypeOf(pinger)
	if resp, err := pinger.Ping(context.Background(), &emptypb.Empty{}); err != nil {
		logger.Error().Msgf("cannot ping %s: %v", pingerType.String(), err)
	} else {
		if resp.GetStatus() != genericproto.ResultStatus_OK {
			logger.Error().Msgf("cannot ping %s: %v", pingerType.String(), resp.GetStatus())
		} else {
			logger.Info().Msgf("%s ping response: %s", pingerType.String(), resp.GetMessage())
		}
	}
}
