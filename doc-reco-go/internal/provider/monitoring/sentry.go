package monitoring

import (
	"doc-reco-go/internal/config"
	"fmt"
	"github.com/getsentry/sentry-go"
	"time"
)

func InitializeSentry() error {
	conf := config.Config.Sentry

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              conf.Dsn,
		Environment:      conf.Env,
		Debug:            conf.Debug,
		AttachStacktrace: true,
	})
	if err != nil {
		return fmt.Errorf("sentry initialization error: %s", err)
	}
	defer sentry.Flush(2 * time.Second)
	return nil
}
