package grpc_server

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/swallowarc/tictactoe_battle_backend/internal/commons/mode"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

type (
	GRPCServer interface {
		RunGRPCServer(ctx context.Context)
	}

	ControllerRegister interface {
		Register(grpcServer grpc.ServiceRegistrar)
	}

	InitFunc   func()
	CloserFunc func()

	grpcServer struct {
		logger             *zap.Logger
		port               string
		mode               mode.Mode
		controllerRegister ControllerRegister
		initFunction       InitFunc
		closerFunction     CloserFunc
	}
)

var (
	catchSignals = []os.Signal{
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	}
)

func NewGRPCServer(
	logger *zap.Logger,
	port string, mode mode.Mode,
	controllerRegister ControllerRegister,
	initFunction InitFunc,
	closerFunction CloserFunc,
) GRPCServer {
	return &grpcServer{
		logger:             logger,
		port:               port,
		mode:               mode,
		controllerRegister: controllerRegister,
		initFunction:       initFunction,
		closerFunction:     closerFunction,
	}
}

func (g *grpcServer) RunGRPCServer(ctx context.Context) {
	g.logger.Info("Initialize gRPC Server ...")
	g.initFunction()

	g.logger.Info("Startup gRPC Server ...")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", g.port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := g.newServer()

	g.logger.Info(fmt.Sprintf("Startup using port : %s", g.port))
	go func() {
		if err := server.Serve(lis); err != nil {
			g.logger.Panic("failed to gRPC Server running", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, catchSignals...)

	select {
	case v := <-ctx.Done():
		g.logger.Info("!! Receive context cancel !!", zap.Reflect("message", v))
	case sig := <-quit:
		g.logger.Info("!! Receive signal !!", zap.String("signal", sig.String()))
	}

	wait := time.Duration(5)
	if g.mode == mode.Debug {
		wait = 1
	}
	sdCtx, sdCancel := context.WithTimeout(ctx, wait*time.Second)
	defer sdCancel()

	<-sdCtx.Done()

	g.logger.Info("Closing gRPC Server ...")
	server.GracefulStop()
	g.closerFunction()

	g.logger.Info("Shutdown gRPC Server")
}

func (g *grpcServer) newServer() *grpc.Server {
	var (
		zapOpts = []grpc_zap.Option{
			grpc_zap.WithDurationField(func(duration time.Duration) zapcore.Field {
				return zap.Int64("grpc.time_ns", duration.Nanoseconds())
			}),
			// Shared options for the logger, with a custom gRPC code to log level function.
			//grpc_zap.WithLevels(func(code codes.Code) zapcore.Level {
			//	switch code {
			//	case codes.OK, codes.Canceled:
			//		return zapcore.DebugLevel
			//	// case ...
			//	default:
			//		return zapcore.ErrorLevel
			//	}
			//}),
		}

		//kaep = keepalive.EnforcementPolicy{
		//	MinTime:             10 * time.Second,
		//	PermitWithoutStream: true,
		//}

		kasp = keepalive.ServerParameters{
			MaxConnectionIdle: 600 * time.Second,
		}
	)

	// TODO: add interceptors(ex: auth)

	// Make sure that log statements internal to gRPC library are logged using the zapLogger as well.
	grpc_zap.ReplaceGrpcLoggerV2(g.logger)
	// Create a server, make sure we put the grpc_ctxtags context before everything else.
	server := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_zap.UnaryServerInterceptor(g.logger, zapOpts...),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_zap.StreamServerInterceptor(g.logger, zapOpts...),
		),
		//grpc.KeepaliveEnforcementPolicy(kaep),
		grpc.KeepaliveParams(kasp),
	)

	g.controllerRegister.Register(server)

	if g.mode == mode.Debug {
		reflection.Register(server)
	}
	return server
}
