package resolver

import (
	"context"
	"crypto/tls"
	"emperror.dev/errors"
	pb "github.com/je4/miniresolver/v2/pkg/miniresolverproto"
	trusthelper "github.com/je4/trustutil/v2/pkg/grpchelper"
	"github.com/je4/utils/v2/pkg/zLogger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
	"strconv"
	"sync"
	"time"
)

func newServer(addr string, domains []string, tlsConfig *tls.Config, resolver pb.MiniResolverClient, single bool, logger zLogger.ZLogger, opts ...grpc.ServerOption) (*Server, error) {
	if tlsConfig == nil {
		return nil, errors.New("no tls configuration")
	}
	listenConfig := &net.ListenConfig{
		Control:   nil,
		KeepAlive: 0,
	}
	lis, err := listenConfig.Listen(context.Background(), "tcp", addr)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot listen on %s", addr)
	}
	addr = lis.Addr().String()
	logger.Info().Msgf("listening on %s", addr)
	l2 := logger.With().Str("addr", addr).Logger()
	logger = &l2
	interceptor := trusthelper.NewInterceptor(domains, logger)

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
		domains:      domains,
		single:       single,
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
	domains      []string
	single       bool
}

func (s *Server) GetAddr() string {
	return s.addr
}

func (s *Server) Startup() {
	s.waitShutdown.Add(2)
	go func() {
		defer s.waitShutdown.Done()
		s.logger.Info().Msg("starting server")
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
		singlestr := ""
		if s.single {
			singlestr = "single "
		}
		var endLoop = false
		for endLoop == false {
			var waitSeconds int64 = 10
			for name, _ := range si {
				s.logger.Info().Msgf("registering %sservice %v.%s at %s ", singlestr, s.domains, name, s.addr)
				_, port, err := net.SplitHostPort(s.addr)
				if err != nil {
					s.logger.Error().Err(err).Msgf("cannot split host port of '%s'", s.addr)
					continue
				}
				portInt, err := strconv.Atoi(port)
				if err != nil {
					s.logger.Error().Err(err).Msgf("cannot convert port '%s' to int", port)
					continue
				}
				if resp, err := s.resolver.AddService(context.Background(), &pb.ServiceData{
					Service: name,
					Port:    uint32(portInt),
					Domains: s.domains,
					Single:  s.single,
				}); err != nil {
					s.logger.Error().Err(err).Msg("cannot register service")
				} else {
					waitSeconds = resp.GetNextCallWait()
					s.logger.Info().Msgf("%sservice registered: %v", singlestr, resp.GetResponse().GetMessage())
				}
			}
			if waitSeconds == 0 {
				waitSeconds = 5 * 60
			}
			s.logger.Debug().Msgf("waiting %d seconds for refreshing service", waitSeconds)
			select {
			case <-s.done:
				endLoop = true
				s.logger.Info().Msg("ending resolver refresh loop")
				break
			case <-time.After(time.Duration(waitSeconds) * time.Second):
			}
		}
		for name, _ := range si {
			s.logger.Info().Msgf("unregistering %sservice %v.%s at %s", singlestr, s.domains, name, s.addr)
			_, port, err := net.SplitHostPort(s.addr)
			if err != nil {
				s.logger.Error().Err(err).Msgf("cannot split host port of '%s'", s.addr)
				continue
			}
			portInt, err := strconv.Atoi(port)
			if err != nil {
				s.logger.Error().Err(err).Msgf("cannot convert port '%s' to int", port)
				continue
			}
			if resp, err := s.resolver.RemoveService(context.Background(), &pb.ServiceData{
				Service: name,
				Port:    uint32(portInt),
				Domains: s.domains,
			}); err != nil {
				s.logger.Error().Err(err).Msg("cannot unregister service")
			} else {
				s.logger.Info().Msgf("%sservice unregistered: %v", singlestr, resp.Message)
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
