package main

import (
	"fmt"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
	"strings"
)

func main() {
	proxywasm.SetVMContext(&vmContext{})
}

type vmContext struct {
	types.DefaultVMContext
}

func (*vmContext) NewPluginContext(id uint32) types.PluginContext {
	return &pluginContext{}
}

type pluginContext struct {
	types.DefaultPluginContext
}

func (*pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	return &httpHeaders{contextID: contextID}
}

type httpHeaders struct {
	types.DefaultHttpContext
	contextID uint32
}

func (h *httpHeaders) OnHttpRequestHeaders(_ int, _ bool) types.Action {
	rawHeaderValue, err := proxywasm.GetHttpRequestHeader("authorization")
	if err != nil {
		return sendUnauthorized(fmt.Sprintf("can't get authorization header: %v", err))
	}

	proxywasm.LogInfof("Authorization header: %s", rawHeaderValue)
	token := strings.TrimPrefix(rawHeaderValue, "Bearer ")
	if token != "foobarbaz" {
		return sendUnauthorized("incorrect token")
	}

	return types.ActionContinue
}

func sendUnauthorized(why string) types.Action {
	proxywasm.LogInfof("sendUnauthorized: %s", why)
	hs := [][2]string{
		{"foo", "bar"},
	}
	_ = proxywasm.SendHttpResponse(401, hs, []byte("unauthorized"))
	return types.ActionPause
}
