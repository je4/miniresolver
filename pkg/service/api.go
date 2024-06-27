package service

import (
	"context"
	"emperror.dev/errors"
	"fmt"
	"github.com/elazarl/goproxy"
	pbgeneric "github.com/je4/genericproto/v2/pkg/generic/proto"
	pb "github.com/je4/miniresolver/v2/pkg/miniresolverproto"
	"github.com/je4/utils/v2/pkg/zLogger"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"io"
	"net"
	"net/http"
	"time"
)

func NewMiniResolver(bufferSize int, serviceExpiration time.Duration, proxy string, logger zLogger.ZLogger) *miniResolver {
	_logger := logger.With().Str("rpcService", "miniResolver").Logger()
	return &miniResolver{
		logger:            &_logger,
		services:          newCache(serviceExpiration, &_logger),
		serviceExpiration: serviceExpiration,
		proxyAddr:         proxy,
	}
}

type miniResolver struct {
	pb.UnimplementedMiniResolverServer
	logger            zLogger.ZLogger
	services          *cache
	serviceExpiration time.Duration
	proxyAddr         string
	proxyServer       *http.Server
}

func (d *miniResolver) StartProxy() error {
	handler := goproxy.NewProxyHttpServer()
	d.proxyServer = &http.Server{
		Addr:    d.proxyAddr,
		Handler: handler,
	}
	handler.Verbose = true
	handler.OnRequest().HijackConnect(func(req *http.Request, client net.Conn, ctx *goproxy.ProxyCtx) {
		addr, _ := d.services.getService(req.URL.Host)
		if addr == "" {
			d.logger.Debug().Str("proxy", "HijackConnect()").Msgf("service '%s' not found", req.URL.Host)
			client.Write([]byte("HTTP/1.1 404 Service not found\r\n\r\n"))
			return
		}
		d.logger.Debug().Str("proxy", "HijackConnect()").Msgf("proxy connect to %s", req.URL.Host)
		defer func() {
			if e := recover(); e != nil {
				d.logger.Debug().Str("proxy", "HijackConnect()").Msgf("error connecting to remote: %v", e)
				client.Write([]byte("HTTP/1.1 500 Cannot reach destination\r\n\r\n"))
			}
			client.Close()
		}()
		//clientBuf := bufio.NewReadWriter(bufio.NewReader(client), bufio.NewWriter(client))
		remote, err := net.Dial("tcp", addr)
		if err != nil {
			ctx.Logf("error connecting to remote: %v", err)
			client.Write([]byte("HTTP/1.1 500 Cannot reach destination\r\n\r\n"))
			return
		}
		client.Write([]byte("HTTP/1.1 200 Ok\r\n\r\n"))

		//remoteBuf := bufio.NewReadWriter(bufio.NewReader(remote), bufio.NewWriter(remote))

		done := make(chan bool)

		go func() {
			if _, err := io.Copy(client, remote); err != nil {
				d.logger.Debug().Str("proxy", "HijackConnect()").Msgf("error copying from remote to client: %v", err)
			}
			defer client.Close()
			done <- true
		}()

		go func() {
			if _, err := io.Copy(remote, client); err != nil {
				d.logger.Debug().Str("proxy", "HijackConnect()").Msgf("error copying from client to remote: %v", err)
			}
			defer remote.Close()
			done <- true
		}()

		<-done
		<-done

		/*
			for {
				req, err := http.ReadRequest(clientBuf.Reader)
				if err != nil {
					ctx.Logf("error reading request: %v", err)
					return
				}
				if err := req.Write(remoteBuf); err != nil {
					ctx.Logf("error writing request: %v", err)
					return
				}
				if err := remoteBuf.Flush(); err != nil {
					ctx.Logf("error flushing request: %v", err)
					return
				}
				resp, err := http.ReadResponse(remoteBuf.Reader, req)
				if err != nil {
					ctx.Logf("error reading response: %v", err)
					return
				}
				if err := resp.Write(clientBuf.Writer); err != nil {
					ctx.Logf("error writing response: %v", err)
					return
				}
				if err := clientBuf.Flush(); err != nil {
					ctx.Logf("error flushing response: %v", err)
					return
				}
			}

		*/
		d.logger.Debug().Str("proxy", "HijackConnect()").Msgf("proxy connect to %s done", req.URL.Host)
	})
	go func() {
		d.logger.Debug().Str("proxy", "HijackConnect()").Msgf("starting proxy on %s", d.proxyAddr)
		if err := d.proxyServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			d.logger.Error().Str("proxy", "HijackConnect()").Msgf("cannot start proxy: %v", err)
		}
	}()
	return nil
}

func (d *miniResolver) StopProxy() error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	d.logger.Debug().Msgf("shutting down proxy")
	if err := d.proxyServer.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "cannot shutdown proxy")
	}
	return nil
}

func (d *miniResolver) Ping(context.Context, *emptypb.Empty) (*pbgeneric.DefaultResponse, error) {
	return &pbgeneric.DefaultResponse{
		Status:  pbgeneric.ResultStatus_OK,
		Message: "pong",
		Data:    nil,
	}, nil
}

func (d *miniResolver) AddService(ctx context.Context, data *pb.ServiceData) (*pb.ResolverDefaultResponse, error) {
	d.logger.Debug().Msgf("add service '%v.%s' - '%s:%d'", data.GetDomains(), data.GetService(), data.GetHost(), data.GetPort())

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
	waitSeconds := int64((d.serviceExpiration.Seconds() * 2.0) / 3.0)
	d.services.addService(data.GetService(), address, data.GetDomains(), data.GetSingle())
	d.logger.Debug().Msgf("service '%s' - '%s' added", data.Service, address)
	return &pb.ResolverDefaultResponse{
		Response: &pbgeneric.DefaultResponse{
			Status:  pbgeneric.ResultStatus_OK,
			Message: fmt.Sprintf("service '%v.%s' - '%s' added", data.Domains, data.Service, address),
		},
		NextCallWait: waitSeconds,
	}, nil
}

func (d *miniResolver) RemoveService(ctx context.Context, data *pb.ServiceData) (*pbgeneric.DefaultResponse, error) {
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
	d.services.removeService(data.Service, address, data.Domains)
	d.logger.Debug().Msgf("service '%s' - '%s' removed", data.Service, address)
	return &pbgeneric.DefaultResponse{
		Status:  pbgeneric.ResultStatus_OK,
		Message: fmt.Sprintf("service '%s' - '%s' removed", data.Service, address),
	}, nil
}

func (d *miniResolver) ResolveServices(ctx context.Context, data *wrapperspb.StringValue) (*pb.ServicesResponse, error) {
	d.logger.Debug().Msgf("resolve services '%s'", data.Value)
	addrs, ncw := d.services.getServices(data.Value)
	return &pb.ServicesResponse{
		Addrs:        addrs,
		NextCallWait: int64(ncw.Seconds()),
	}, nil
}

func (d *miniResolver) ResolveService(ctx context.Context, data *wrapperspb.StringValue) (*pb.ServiceResponse, error) {
	d.logger.Debug().Msgf("resolve service '%s'", data.Value)
	addr, ncw := d.services.getService(data.Value)
	if addr == "" {
		return nil, fmt.Errorf("service '%s' not found", data.Value)
	}
	return &pb.ServiceResponse{
		Addr:         addr,
		NextCallWait: int64(ncw.Seconds()),
	}, nil
}
