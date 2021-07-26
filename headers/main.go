package main

import (
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func main() {
	proxywasm.SetVMContext(&vmContext{})
}

type vmContext struct {
	// embed noop implementation of vmcontext
	// so we don't need to implement all the methods
	types.DefaultVMContext
}

func (*vmContext) NewPluginContext(id uint32) types.PluginContext {
	return &pluginContext{}
}

type pluginContext struct {
	// embed noop implementation of vmcontext
	// so we don't need to implement all the methods
	types.DefaultPluginContext
}

func (*pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	return &httpHeaders{contextID: contextID}
}

type httpHeaders struct {
	types.DefaultHttpContext
	contextID uint32
}

func (h *httpHeaders) OnHttpRequestHeaders(num int, endOfStream bool) types.Action {

	proxywasm.LogInfof("OnHttpRequestHeaders(num: %d, endOfStream: %t)", num, endOfStream)

	if err := proxywasm.AddHttpRequestHeader("x-wasm", "hello world"); err != nil {
		proxywasm.LogCriticalf("unable to add request header: %v", err)
	}

	if err := proxywasm.RemoveHttpRequestHeader("user-agent"); err != nil {
		proxywasm.LogCriticalf("unable to delete request header: %v", err)
	}

	header, err := proxywasm.GetHttpRequestHeader("Accept")
	if err != nil {
		proxywasm.LogCriticalf("unable to get accept request header: %v", err)
	}
	proxywasm.LogInfof("accept: %s", header)

	if err := proxywasm.ReplaceHttpRequestHeader("Accept", "text/plain"); err != nil {
		proxywasm.LogCriticalf("unable to replace accept request header: %v", err)
	}
	header, err = proxywasm.GetHttpRequestHeader("Accept")
	if err != nil {
		proxywasm.LogCriticalf("unable to get accept request header: %v", err)
	}
	proxywasm.LogInfof("accept: %s", header)

	hs, err := proxywasm.GetHttpRequestHeaders()
	if err != nil {
		proxywasm.LogCriticalf("unable to get request headers: %v", err)
	}

	proxywasm.LogInfof("numHeaders: %d", len(hs))
	for _, header := range hs {
		proxywasm.LogInfof("header: %#v", header)
	}

	var headers [][2]string
	proxywasm.SendHttpResponse(401, headers, []byte("Unauthorized"))

	headersSlice := [][2]string{
		{"foo", "bar"},
		{"bar", "baz"},
	}
	body := []byte("hello world")
	timeoutMs := uint32(1000)

	callback := func(numHeaders, bodySize, numTrailers int) {
		// During callBack is called, "GetHttpCallResponseHeaders", "GetHttpCallResponseBody", "GetHttpCallResponseTrailers"
		// calls are available for accessing the response information.
		proxywasm.LogInfof("numHeaders: %d, bodySize: %d, numTrailers: %d", numHeaders, bodySize, numTrailers)
	}

	calloutID, err := proxywasm.DispatchHttpCall(
		"https://httpbin.org/xml",
		headersSlice,
		body,
		nil,
		timeoutMs,
		callback,
		)
	if err != nil {
		proxywasm.LogCriticalf("unable to DispatchHttpCall: %v", err)
	}
	proxywasm.LogInfof("calloutID: %d", calloutID)

//	HTTP calls
//	proxy_dispatch_http_call
//
//params:
//	i32 (const char*) upstream_name_data
//	i32 (size_t) upstream_name_size
//	i32 (const char*) headers_map_data
//	i32 (size_t) headers_map_size
//	i32 (const char*) body_data
//	i32 (size_t) body_size
//	i32 (const char*) trailers_map_data
//	i32 (size_t) trailers_map_size
//	i32 (uint32_t) timeout_milliseconds
//	i32 (uint32_t*) return_callout_id
//returns:
//	i32 (proxy_result_t) call_result
//
//	Dispatch a HTTP call to upstream (upstream_name_data, upstream_name_size). Once the response is returned to the host, proxy_on_http_call_response will be called with a unique call identifier (return_callout_id).
//

	return types.ActionPause
}

func (h *httpHeaders) OnHttpResponseHeaders(num int, endOfStream bool) types.Action {

	proxywasm.LogInfof("OnHttpResponseHeaders(num: %d, endOfStream: %t)", num, endOfStream)

	if err := proxywasm.AddHttpResponseHeader("x-asoorm-added", "woo hoo"); err != nil {
		proxywasm.LogInfof("unable to add http response header: %v", err)
	}

	//svr, err := proxywasm.GetHttpResponseHeader("server")
	//if err != nil {
	//	proxywasm.LogCriticalf("unable to get server response header: %v", err)
	//}
	//
	//proxywasm.LogInfof("server: %s", svr)

	//if err := proxywasm.ReplaceHttpResponseHeader("server", "asoorm"); err != nil {
	//	proxywasm.LogCriticalf("unable to replace server response header: %v", err)
	//}

	//svr, err = proxywasm.GetHttpResponseHeader("server")
	//if err != nil {
	//	proxywasm.LogCriticalf("unable to get server response header: %v", err)
	//}
	//
	//proxywasm.LogInfof("server: %s", svr)
	//
	//err = proxywasm.RemoveHttpResponseHeader("Access-Control-Allow-Origin")
	//if err != nil {
	//	proxywasm.LogCriticalf("unable to remove Access-Control-Allow-Origin: %v", err)
	//}

	return types.ActionContinue
}
