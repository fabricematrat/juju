// Copyright 2012, 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package base

// APICaller is implemented by the client-facing State object.
type APICaller interface {
	// Call makes a call to the API server with the given object type,
	// id, request and parameters. The response is filled in with the
	// call's result if the call is successful.
	APICall(objType string, version int, id, request string, params, response interface{}) error

	// BestFacadeVersion returns the newest version of 'objType' that this
	// client can use with the current API server.
	BestFacadeVersion(facade string) int
}

// FacadeCaller is a wrapper around the common paradigm that a given client
// just wants to make calls on a facade using the best known version of the API.
type FacadeCaller interface {
	// CallFacade will place a request against the API using the requested
	// Facade and the best version that the API server supports that is
	// also known to the client.
	CallFacade(request string, params, response interface{}) error

	// BestAPIVersion returns the API version that we were able to
	// determine is supported by both the client and the API Server
	BestAPIVersion() int

	// RawAPICaller returns the wrapped APICaller. This can be used if you need
	// to switch what Facade you are calling (such as Facades that return
	// Watchers and then need to use the Watcher facade)
	RawAPICaller() APICaller
}

type facadeCaller struct {
	facade string
	caller APICaller
}

func (fc facadeCaller) CallFacade(request string, params, response interface{}) error {
	return fc.caller.APICall(
		fc.facade, fc.caller.BestFacadeVersion(fc.facade), "",
		request, params, response)
}

func (fc facadeCaller) BestAPIVersion() int {
	// Note: If we decide we want to cache the best version rather than
	// computing it from a list each time, this is a clear location to do
	// that caching.
	return fc.caller.BestFacadeVersion(fc.facade)
}

func (fc facadeCaller) RawAPICaller() APICaller {
	return fc.caller
}

// GetFacadeCaller wraps a APICaller for a given Facade
func GetFacadeCaller(caller APICaller, facade string) FacadeCaller {
	return facadeCaller{facade: facade, caller: caller}
}
