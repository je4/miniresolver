package resolver

import (
	"crypto/tls"
	"emperror.dev/errors"
	"fmt"
	"github.com/je4/miniresolver/v2/pkg/builder"
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

/*
NewMiniresolverClient creates a new miniresolver client
serverAddr: network address of the miniresolver server can be empty for clientMap only
clientMap: map of service names to network addresses can be nil
clientTLSConfig: tls configuration for the client can be nil for server only or insecure
serverTLSConfig: tls configuration for the server can be nil for client only or insecure
resolverTimeout: timeout for resolver recall
resolverNotFoundTimeout: timeout for resolver not found for recall if resolver not found
logger: logger
opts: additional grpc dial options
*/
func NewMiniresolverClient(serverAddr string, clientMap map[string]string, clientTLSConfig, serverTLSConfig *tls.Config, resolverTimeout, resolverNotFoundTimeout time.Duration, logger zLogger.ZLogger, dialOpts ...grpc.DialOption) (*MiniResolver, error) {
	var err error
	if clientMap == nil {
		clientMap = map[string]string{}
	}
	if dialOpts == nil {
		dialOpts = []grpc.DialOption{}
	}
	res := &MiniResolver{
		clientMap:       clientMap,
		clientTLSConfig: clientTLSConfig,
		serverTLSConfig: serverTLSConfig,
		dialOpts:        dialOpts,
		serverOpts:      []grpc.ServerOption{},
		logger:          logger,
	}
	if serverAddr != "" {
		res.MiniResolverClient, res.conn, err = newClient[pb.MiniResolverClient](pb.NewMiniResolverClient, serverAddr, clientTLSConfig, dialOpts...)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot create client for %s", serverAddr)
		}
		builder.RegisterResolver(res.MiniResolverClient, resolverTimeout, resolverNotFoundTimeout, logger)
	}

	return res, nil
}

type MiniResolver struct {
	pb.MiniResolverClient
	conn            io.Closer
	clientCloser    []io.Closer
	clientTLSConfig *tls.Config
	dialOpts        []grpc.DialOption
	serverOpts      []grpc.ServerOption
	serverTLSConfig *tls.Config
	clientMap       map[string]string
	logger          zLogger.ZLogger
}

func (c *MiniResolver) SetDialOpts(options ...grpc.DialOption) {
	c.dialOpts = append(c.dialOpts, options...)
}

func (c *MiniResolver) SetServerOpts(options ...grpc.ServerOption) {
	c.serverOpts = append(c.serverOpts, options...)
}

func NewClient[V any](c *MiniResolver, newClientFunc func(conn grpc.ClientConnInterface) V, serviceName string) (V, error) {
	var n V
	var clientAddr string
	if _, ok := c.clientMap[serviceName]; ok {
		clientAddr = c.clientMap[serviceName]
	} else {
		if c.MiniResolverClient != nil {
			clientAddr = fmt.Sprintf("miniresolver:%s", serviceName)
		}
	}
	if clientAddr == "" {
		return n, errors.Errorf("cannot find client address for %s", serviceName)
	}
	client, conn, err := newClient[V](newClientFunc, clientAddr, c.clientTLSConfig, c.dialOpts...)
	if err != nil {
		return n, errors.Wrapf(err, "cannot create client for %s", clientAddr)
	}
	c.clientCloser = append(c.clientCloser, conn)
	return client, nil
}

func (c *MiniResolver) Close() error {
	var errs []error
	for _, closer := range append(c.clientCloser, c.conn) {
		if e := closer.Close(); e != nil {
			errs = append(errs, e)
		}
	}
	return errors.Combine(errs...)
}

func (c *MiniResolver) NewServer(addr string) (*Server, error) {
	if c.MiniResolverClient == nil {
		return nil, errors.Errorf("no miniresolver client")
	}
	server, err := newServer(addr, c.serverTLSConfig, c.MiniResolverClient, c.logger, c.serverOpts...)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot create server for %s", addr)
	}
	return server, nil
}
