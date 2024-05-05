package main

import (
	"flag"
	"os"
	"time"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/joho/godotenv"
	coreService "github.com/nartvt/go-core/service"
	profile "github.com/nartvt/profile-service/api/profile/v1"
	"github.com/nartvt/profile-service/internal/conf"
	"github.com/nartvt/profile-service/internal/nat"
	"github.com/nartvt/profile-service/internal/service"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/env"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"

	logcore "github.com/nartvt/go-core/log"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func initService(logger log.Logger, hs *http.Server, gs *grpc.Server, chartSvc *service.ChartService, profileSvc *service.ProfileService) coreService.Service {
	profile.RegisterChartServiceHTTPServer(hs, chartSvc)
	profile.RegisterChartServiceServer(gs, chartSvc)

	profile.RegisterProfileServiceHTTPServer(hs, profileSvc)
	profile.RegisterProfileServiceServer(gs, profileSvc)
	return coreService.NewService(logger, hs, gs, kratos.ID(id), kratos.Name(Name), kratos.Version(Version))
}

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
	_ = godotenv.Load()
}

func main() {
	flag.Parse()

	c := config.New(
		config.WithSource(
			env.NewSource("IND_"),
			file.NewSource(flagconf),
		),
	)
	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	log.DefaultLogger = log.With(logcore.LogrusConfig(bc.Server.Log),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace_id", tracing.TraceID(),
		"span_id", tracing.SpanID(),
	)

	app, cleanup, err := initApp(bc.Server, bc.Data, log.DefaultLogger)
	if err != nil {
		panic(err)
	}

	initNat(&bc)

	//	go simulator()

	defer closeInfra(cleanup)

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func closeInfra(cleanup func()) {
	cleanup()
	nat.CloseNats()
}

// TODO: Need to enhance here
func initNat(bc *conf.Bootstrap) {
	nat.InitNats(bc.NatServer)
}

func simulator() {
	for {
		log.Infof("PUBLISH MESSAGE: %s", "e29248e4-1c57-45f9-a5a1-d4ac1b198de1")
		nat.PublishNewProfile("ZY0DTL", "e29248e4-1c57-45f9-a5a1-d4ac1b198de1")
		time.Sleep(5 * time.Second)
	}
}
