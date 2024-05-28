package client

import (
	"crypto/tls"
	"emperror.dev/errors"
	"fmt"
	"github.com/je4/miniresolver/v2/pkg/grpchelper"
	pb "github.com/je4/miniresolver/v2/pkg/miniresolverproto"
	"github.com/je4/utils/v2/pkg/zLogger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"time"
)

func newClient[V any](newClientFunc func(conn grpc.ClientConnInterface) V, serverAddr string, tlsConfig *tls.Config, opts ...grpc.DialOption) (V, io.Closer, error) {
	if tlsConfig != nil {
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	conn, err := grpc.NewClient(serverAddr, opts...)

	if err != nil {
		var n V
		return n, nil, errors.Wrapf(err, "cannot connect to %s", serverAddr)
	}
	client := newClientFunc(conn)
	return client, conn, nil
}

type MiniResolverClient struct {
	pb.MiniResolverClient
	conn         io.Closer
	clientCloser []io.Closer
	tlsConfig    *tls.Config
	opts         []grpc.DialOption
}

func NewClient[V any](c *MiniResolverClient, newClientFunc func(conn grpc.ClientConnInterface) V, serviceDesc grpc.ServiceDesc) (V, error) {
	serverAddr := fmt.Sprintf("miniresolver:%s", serviceDesc.ServiceName)
	client, conn, err := newClient[V](newClientFunc, serverAddr, c.tlsConfig, c.opts...)
	if err != nil {
		var n V
		return n, errors.Wrapf(err, "cannot create client for %s", serverAddr)
	}
	c.clientCloser = append(c.clientCloser, conn)
	return client, nil
}

func (c *MiniResolverClient) Close() error {
	var errs []error
	for _, closer := range append(c.clientCloser, c.conn) {
		if e := closer.Close(); e != nil {
			errs = append(errs, e)
		}
	}
	return errors.Combine(errs...)
}

func NewMiniresolverClient(serverAddr string, tlsConfig *tls.Config, resolverTimeout, resolverNotFoundTimeout time.Duration, logger zLogger.ZLogger, opts ...grpc.DialOption) (*MiniResolverClient, error) {
	client, conn, err := newClient[pb.MiniResolverClient](pb.NewMiniResolverClient, serverAddr, tlsConfig, opts...)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot create client for %s", serverAddr)
	}
	grpchelper.RegisterResolver(client, time.Duration(resolverTimeout), time.Duration(resolverNotFoundTimeout), logger)

	return &MiniResolverClient{
		conn:               conn,
		MiniResolverClient: client,
		tlsConfig:          tlsConfig,
		opts:               opts,
	}, nil
}
