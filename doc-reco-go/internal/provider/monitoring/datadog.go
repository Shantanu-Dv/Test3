package monitoring

import (
	"doc-reco-go/internal/config"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
	"net"
)

func InitializeDatadog() error {
	conf := config.Config.Datadog

	addr := net.JoinHostPort(
		conf.AgentHost,
		conf.AgentPort,
	)

	tracer.Start(
		tracer.WithDebugMode(false),
		tracer.WithEnv(conf.Env),
		tracer.WithService(conf.ServiceName),
		tracer.WithServiceVersion(conf.Version),
		tracer.WithAnalytics(true),
		tracer.WithRuntimeMetrics(),
		tracer.WithAgentAddr(addr),
	)

	err := profiler.Start(
		profiler.WithEnv(conf.Env),
		profiler.WithService(conf.ServiceName),
		profiler.WithVersion(conf.Version),
		profiler.WithAgentAddr(addr),
	)
	if err != nil {
		return err
	}

	return nil
}

func StopDdTracer() {
	tracer.Stop()
	profiler.Stop()
}
