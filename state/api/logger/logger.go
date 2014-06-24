// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package logger

import (
	"fmt"

	"github.com/juju/juju/state/api/base"
	"github.com/juju/juju/state/api/params"
	"github.com/juju/juju/state/api/watcher"
)

// State provides access to an logger worker's view of the state.
type State struct {
	caller base.APICaller
}

func (st *State) call(method string, params, result interface{}) error {
	return st.caller.APICall(
		"Logger", st.caller.BestFacadeVersion("Logger"), "",
		method, params, result)
}

// NewState returns a version of the state that provides functionality
// required by the logger worker.
func NewState(caller base.APICaller) *State {
	return &State{caller}
}

// LoggingConfig returns the loggo configuration string for the agent
// specified by agentTag.
func (st *State) LoggingConfig(agentTag string) (string, error) {
	var results params.StringResults
	args := params.Entities{
		Entities: []params.Entity{{Tag: agentTag}},
	}
	err := st.call("LoggingConfig", args, &results)
	if err != nil {
		// TODO: Not directly tested
		return "", err
	}
	if len(results.Results) != 1 {
		// TODO: Not directly tested
		return "", fmt.Errorf("expected 1 result, got %d", len(results.Results))
	}
	result := results.Results[0]
	if err := result.Error; err != nil {
		return "", err
	}
	return result.Result, nil
}

// WatchLoggingConfig returns a notify watcher that looks for changes in the
// logging-config for the agent specifed by agentTag.
func (st *State) WatchLoggingConfig(agentTag string) (watcher.NotifyWatcher, error) {
	var results params.NotifyWatchResults
	args := params.Entities{
		Entities: []params.Entity{{Tag: agentTag}},
	}
	err := st.call("WatchLoggingConfig", args, &results)
	if err != nil {
		// TODO: Not directly tested
		return nil, err
	}
	if len(results.Results) != 1 {
		// TODO: Not directly tested
		return nil, fmt.Errorf("expected 1 result, got %d", len(results.Results))
	}
	result := results.Results[0]
	if result.Error != nil {
		//  TODO: Not directly tested
		return nil, result.Error
	}
	w := watcher.NewNotifyWatcher(st.caller, result)
	return w, nil
}
