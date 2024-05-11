package grpchelper

import (
	"context"
	"crypto/tls"
	"emperror.dev/errors"
	pb "github.com/je4/miniresolver/v2/pkg/miniresolverproto"
	trusthelper "github.com/je4/trustutil/v2/pkg/grpchelper"
	"github.com/je4/trustutil/v2/pkg/tlsutil"
	"github.com/je4/utils/v2/pkg/zLogger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
	"sync"
	"time"
)

func NewServer(addr string, tlsConfig *tls.Config, resolver pb.MiniResolverClient, logger zLogger.ZLogger, opts ...grpc.ServerOption) (*Server, error) {
	listenConfig := &net.ListenConfig{
		Control:   nil,
		KeepAlive: 0,
	}
	lis, err := listenConfig.Listen(context.Background(), "tcp", addr)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot listen on %s", addr)
	}
	interceptor := trusthelper.NewInterceptor(logger)

	if tlsConfig == nil {
		tlsConfig, err = tlsutil.CreateDefaultServerTLSConfig("devServer")
		if err != nil {
			return nil, errors.Wrap(err, "cannot create default server TLS config")
		}
	}

	opts = append(opts, grpc.Creds(credentials.NewTLS(tlsConfig)), grpc.UnaryInterceptor(interceptor.ServerInterceptor))
	grpcServer := grpc.NewServer(opts...)
	server := &Server{
		addr:         addr,
		Server:       grpcServer,
		listener:     lis,
		logger:       logger,
		done:         make(chan bool),
		resolver:     resolver,
		waitShutdown: sync.WaitGroup{},
	}
	return server, nil
}

type Server struct {
	*grpc.Server
	listener     net.Listener
	logger       zLogger.ZLogger
	done         chan bool
	waitShutdown sync.WaitGroup
	resolver     pb.MiniResolverClient
	addr         string
}

func (s *Server) Startup() {

	s.waitShutdown.Add(2)
	go func() {
		defer s.waitShutdown.Done()
		s.logger.Info().Msgf("starting server at %s", s.listener.Addr().String())
		if err := s.Server.Serve(s.listener); err != nil {
			s.logger.Error().Err(err).Msg("cannot serve")
		} else {
			s.logger.Info().Msg("server stopped")
		}
	}()
	time.Sleep(100 * time.Millisecond)
	go func() {
		defer s.waitShutdown.Done()
		si := s.Server.GetServiceInfo()
		for name, _ := range si {
			s.logger.Info().Msgf("registering service %s at %s", name, s.addr)
			if resp, err := s.resolver.AddService(context.Background(), &pb.ServiceData{Service: name, Address: s.addr}); err != nil {
				s.logger.Error().Err(err).Msg("cannot register service")
			} else {
				s.logger.Info().Msgf("service registered: %v", resp.Message)
			}
		}
		<-s.done
		for name, _ := range si {
			s.logger.Info().Msgf("unregistering service %s at %s", name, s.addr)
			if resp, err := s.resolver.RemoveService(context.Background(), &pb.ServiceData{Service: name, Address: s.addr}); err != nil {
				s.logger.Error().Err(err).Msg("cannot unregister service")
			} else {
				s.logger.Info().Msgf("service unregistered: %v", resp.Message)
			}
		}
	}()
}

func (s *Server) Shutdown() error {
	s.done <- true
	s.Server.GracefulStop()
	s.waitShutdown.Wait()
	return errors.Wrap(s.listener.Close(), "cannot close listener")
}
