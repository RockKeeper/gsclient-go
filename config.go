package gsclient

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"runtime"
	"time"
)

const (
	defaultCheckRequestTimeoutSecs = 120
	defaultMaxNumberOfRetries      = 100
	defaultDelayIntervalMilliSecs  = 500
	version                        = "1.0.0"
)

//Config config for client
type Config struct {
	apiURL                  string
	userUUID                string
	apiToken                string
	userAgent               string
	httpClient              *http.Client
	requestCheckTimeoutSecs time.Duration
	delayInterval           time.Duration
	maxNumberOfRetries      int
	logger                  logrus.Logger
}

//NewConfiguration creates a new config
//
//- Parameters:
//		+ apiURL string: base URL of API.
//		+ uuid string: UUID of user.
//		+ token string: API token.
//		+ debugMode bool: true => run client in debug mode.
//		+ requestCheckTimeoutSecs int: Timeout for checking requests (for synchronous feature)
//		+ delayIntervalMilliSecs int: delay between requests when checking request (or retry 5xx, 424 error code)
//		+ maxNumberOfRetries int: number of retries when server returns 5xx, 424 error code.
func NewConfiguration(apiURL string, uuid string, token string, debugMode bool, requestCheckTimeoutSecs,
	delayIntervalMilliSecs, maxNumberOfRetries int) *Config {
	logLevel := logrus.InfoLevel
	if debugMode {
		logLevel = logrus.DebugLevel
	}

	logger := logrus.Logger{
		Out:   os.Stderr,
		Level: logLevel,
		Formatter: &logrus.TextFormatter{
			FullTimestamp: true,
			DisableColors: false,
		},
	}

	if requestCheckTimeoutSecs == 0 {
		requestCheckTimeoutSecs = defaultCheckRequestTimeoutSecs
	}
	if delayIntervalMilliSecs == 0 {
		delayIntervalMilliSecs = defaultDelayIntervalMilliSecs
	}
	if maxNumberOfRetries == 0 {
		maxNumberOfRetries = defaultMaxNumberOfRetries
	}
	cfg := &Config{
		apiURL:                  apiURL,
		userUUID:                uuid,
		apiToken:                token,
		userAgent:               "gsclient-go/" + version + " (" + runtime.GOOS + ")",
		httpClient:              http.DefaultClient,
		logger:                  logger,
		requestCheckTimeoutSecs: time.Duration(requestCheckTimeoutSecs) * time.Second,
		delayInterval:           time.Duration(delayIntervalMilliSecs) * time.Millisecond,
		maxNumberOfRetries:      maxNumberOfRetries,
	}
	return cfg
}
