package resolver

import (
	"context"
	"crypto/tls"
	"emperror.dev/errors"
	"fmt"
	pb "github.com/je4/miniresolver/v2/pkg/miniresolverproto"
	"github.com/je4/utils/v2/pkg/zLogger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"
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
		watchServices:   map[string]chan<- bool{},
		watchLock:       sync.Mutex{},
		clientMap:       clientMap,
		clientTLSConfig: clientTLSConfig,
		serverTLSConfig: serverTLSConfig,
		dialOpts:        dialOpts,
		serverOpts:      []grpc.ServerOption{},
		logger:          logger,
	}
	//res.SetServerOpts(grpc.ChainUnaryInterceptor(res.unaryServerInterceptor), grpc.ChainStreamInterceptor(res.streamServerInterceptor))
	res.SetDialOpts(
		grpc.WithUnaryInterceptor(res.getUnaryClientInterceptor()),
		grpc.WithStreamInterceptor(res.getStreamClientInterceptor()),
	)
	if serverAddr != "" {
		res.MiniResolverClient, res.conn, err = newClient[pb.MiniResolverClient](pb.NewMiniResolverClient, serverAddr, clientTLSConfig, res.dialOpts...)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot create client for %s", serverAddr)
		}
		RegisterResolver(res, resolverTimeout, resolverNotFoundTimeout, logger)
	}

	return res, nil
}

type MiniResolver struct {
	pb.MiniResolverClient
	conn            io.Closer
	watchLock       sync.Mutex
	watchServices   map[string]chan<- bool
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

var domainRegexp = regexp.MustCompile(`^miniresolver:([a-zA-Z0-9-]+)\.([a-zA-Z0-9-]+)\.([a-zA-Z0-9-]+)`)

func (c *MiniResolver) getStreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		var domain string
		target := cc.Target()
		if matches := domainRegexp.FindStringSubmatch(target); matches != nil {
			domain = matches[1]
		}
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}
		d := md.Get("domain")
		if len(d) == 0 {
			md.Set("domain", domain)
		} else {
			domain = d[0]
		}
		ctx = metadata.NewOutgoingContext(ctx, md)
		start := time.Now()
		clientStream, err := streamer(ctx, desc, cc, method, opts...)
		end := time.Now()
		c.logger.Debug().Str("domain", domain).Str("target", target).Str("method", method).Dur("duration", end.Sub(start)).Err(err)
		if err != nil {
			if stat, ok := status.FromError(err); ok {
				if stat.Code() == codes.Unavailable {
					c.RefreshResolver(cc.Target())
				}
			}
			return nil, errors.Wrapf(err, "RPC: %s %s :: %s", target, method, domain)
		}
		return clientStream, nil
	}
}

func (c *MiniResolver) getUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		start := time.Now()
		var domain string
		var target = cc.Target()
		if matches := domainRegexp.FindStringSubmatch(target); matches != nil {
			domain = matches[1]
		}
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}
		d := md.Get("domain")
		if len(d) == 0 {
			md.Set("domain", domain)
		} else {
			domain = d[0]
		}
		ctx = metadata.NewOutgoingContext(ctx, md)
		err := invoker(ctx, method, req, reply, cc, opts...)
		end := time.Now()
		c.logger.Debug().Str("domain", domain).Str("target", target).Str("method", method).Dur("duration", end.Sub(start)).Err(err)
		if err != nil {
			if status, ok := status.FromError(err); ok {
				if status.Code() == codes.Unavailable {
					c.RefreshResolver(target)
				}
			}
			return errors.Wrapf(err, "RPC: %s/%s :: %s", target, method, domain)
		}
		return nil
	}
}

func NewClients[V any](c *MiniResolver, newClientFunc func(conn grpc.ClientConnInterface) V, serviceName string, domains []string) (map[string]V, error) {
	var result = map[string]V{}
	for _, domain := range domains {
		client, err := NewClient[V](c, newClientFunc, serviceName, domain)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot create client for %s.%s", domain, serviceName)
		}
		result[domain] = client
	}
	return result, nil
}

func NewClient[V any](c *MiniResolver, newClientFunc func(conn grpc.ClientConnInterface) V, serviceName, domain string) (V, error) {
	client, closer, err := NewClientCloser(c, newClientFunc, serviceName, domain)
	if err != nil {
		return client, errors.Wrapf(err, "cannot create client for %s.%s", domain, serviceName)
	}
	c.clientCloser = append(c.clientCloser, closer)
	return client, nil
}

func NewClientCloser[V any](c *MiniResolver, newClientFunc func(conn grpc.ClientConnInterface) V, serviceName, domain string) (V, io.Closer, error) {
	var n V
	var clientAddr string

	if domain != "" {
		serviceName = domain + "." + serviceName
	}

	if _, ok := c.clientMap[serviceName]; ok {
		clientAddr = c.clientMap[serviceName]
	} else {
		if os.Getenv("HTTPS_PROXY") != "" {
			clientAddr = "passthrough:///" + serviceName
		} else {
			if c.MiniResolverClient != nil && !strings.Contains(serviceName, ":") {
				clientAddr = fmt.Sprintf("miniresolver:%s", serviceName)
			}
		}
	}
	if clientAddr == "" {
		return n, nil, errors.Errorf("cannot find client address for %s", serviceName)
	}
	client, conn, err := newClient[V](newClientFunc, clientAddr, c.clientTLSConfig, c.dialOpts...)
	if err != nil {
		return n, nil, errors.Wrapf(err, "cannot create client for %s", clientAddr)
	}
	return client, conn, nil
}

func (c *MiniResolver) Close() error {
	var errs []error
	for _, closer := range c.clientCloser {
		if e := closer.Close(); e != nil {
			errs = append(errs, e)
		}
	}
	if c.conn != nil {
		if e := c.conn.Close(); e != nil {
			errs = append(errs, e)
		}
	}
	return errors.Combine(errs...)
}

func (c *MiniResolver) NewServer(addr string, domains []string, single bool) (*Server, error) {
	if c.MiniResolverClient == nil {
		return nil, errors.Errorf("no miniresolver client")
	}
	server, err := newServer(addr, domains, c.serverTLSConfig, c.MiniResolverClient, single, c.logger, c.serverOpts...)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot create server for %s", addr)
	}
	return server, nil
}

func (c *MiniResolver) WatchService(target string, reloadChannel chan<- bool) {
	c.logger.Debug().Msgf("watch service %s", target)
	c.watchLock.Lock()
	defer c.watchLock.Unlock()
	c.watchServices[target] = reloadChannel
}

func (c *MiniResolver) UnwatchService(target string) {
	c.logger.Debug().Msgf("unwatch service %s", target)
	c.watchLock.Lock()
	defer c.watchLock.Unlock()
	delete(c.watchServices, target)
}

func (c *MiniResolver) RefreshResolver(target string) {
	c.watchLock.Lock()
	ch, ok := c.watchServices[target]
	c.watchLock.Unlock()
	if ok {
		select {
		case ch <- true:
			c.logger.Debug().Msgf("refresh service %s", target)
		case <-time.After(2 * time.Second):
			c.logger.Error().Msgf("cannot refresh resolver for %s", target)
			c.logger.Debug().Msgf("timeout refresh service %s", target)
		}
	} else {
		c.logger.Debug().Msgf("service %s not in watch map", target)
	}
}
