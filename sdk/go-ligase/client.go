package ligase

import (
	"net/http"
	"sync"
	"time"

	"github.com/finogeeks/ligase/sdk/go-ligase/logger"
)

var (
	once      sync.Once
	ligaseURL string
	httpCli   *http.Client
)

const (
	R0_PREFIX = "/_matrix/client/r0"
)

type ClientOpts struct {
	URL                 string
	MaxIdleConns        int
	MaxIdleConnsPerHost int
	MaxConnsPerHost     int
	IdleConnTimeout     time.Duration

	LogTime   bool
	LogFile   bool
	LogDebug  bool
	LogTrace  bool
	LogColors bool
	LogPID    bool
}

// Init 初始化
func Init(opts ClientOpts) {
	ligaseURL = opts.URL
	once.Do(func() {
		logger.InitLogger(opts.LogTime, opts.LogFile, opts.LogDebug, opts.LogTrace, opts.LogColors, opts.LogPID)
		httpCli = &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:        opts.MaxIdleConns,
				MaxIdleConnsPerHost: opts.MaxIdleConnsPerHost,
				MaxConnsPerHost:     opts.MaxConnsPerHost,
				IdleConnTimeout:     opts.IdleConnTimeout,
			},
		}
	})
}
