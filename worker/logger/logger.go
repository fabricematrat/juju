// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package logger

import (
	"github.com/juju/errors"
	"github.com/juju/loggo"
	worker "gopkg.in/juju/worker.v1"

	"github.com/juju/juju/agent"
	"github.com/juju/juju/api/logger"
	"github.com/juju/juju/watcher"
)

var log = loggo.GetLogger("juju.worker.logger")

// Logger is responsible for updating the loggo configuration when the
// environment watcher tells the agent that the value has changed.
type Logger struct {
	api         *logger.State
	agentConfig agent.Config
	lastConfig  string
}

// NewLogger returns a worker.Worker that uses the notify watcher returned
// from the setup.
func NewLogger(api *logger.State, agentConfig agent.Config) (worker.Worker, error) {
	logger := &Logger{
		api:         api,
		agentConfig: agentConfig,
		lastConfig:  loggo.LoggerInfo(),
	}
	log.Debugf("initial log config: %q", logger.lastConfig)
	w, err := watcher.NewNotifyWorker(watcher.NotifyConfig{
		Handler: logger,
	})
	if err != nil {
		return nil, errors.Trace(err)
	}
	return w, nil
}

func (logger *Logger) setLogging() {
	loggingConfig, err := logger.api.LoggingConfig(logger.agentConfig.Tag())
	if err != nil {
		log.Errorf("%v", err)
	} else {
		if loggingConfig != logger.lastConfig {
			log.Debugf("reconfiguring logging from %q to %q", logger.lastConfig, loggingConfig)
			loggo.DefaultContext().ResetLoggerLevels()
			if err := loggo.ConfigureLoggers(loggingConfig); err != nil {
				// This shouldn't occur as the loggingConfig should be
				// validated by the original Config before it gets here.
				log.Warningf("configure loggers failed: %v", err)
				// Try to reset to what we had before
				loggo.ConfigureLoggers(logger.lastConfig)
			}
			logger.lastConfig = loggingConfig
		}
	}
}

func (logger *Logger) SetUp() (watcher.NotifyWatcher, error) {
	log.Debugf("logger setup")
	// We need to set this up initially as the NotifyWorker sucks up the first
	// event.
	logger.setLogging()
	return logger.api.WatchLoggingConfig(logger.agentConfig.Tag())
}

func (logger *Logger) Handle(_ <-chan struct{}) error {
	logger.setLogging()
	return nil
}

func (logger *Logger) TearDown() error {
	// Nothing to cleanup, only state is the watcher
	return nil
}
