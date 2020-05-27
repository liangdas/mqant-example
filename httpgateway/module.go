/**
一定要记得在confin.json配置这个模块的参数,否则无法使用
*/
package httpgateway

import (
	"github.com/liangdas/mqant/conf"
	"github.com/liangdas/mqant/httpgateway"
	"github.com/liangdas/mqant/httpgateway/proto"
	"github.com/liangdas/mqant/log"
	"github.com/liangdas/mqant/module"
	"github.com/liangdas/mqant/module/base"
	"github.com/liangdas/mqant/rpc"
	"github.com/liangdas/mqant/rpc/pb"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
)

var Module = func() module.Module {
	this := new(httpgate)
	return this
}

type httpgate struct {
	basemodule.BaseModule
}

func (self *httpgate) GetType() string {
	//很关键,需要与配置文件中的Module配置对应
	return "httpgate"
}
func (self *httpgate) Version() string {
	//可以在监控时了解代码版本
	return "1.0.0"
}
func (self *httpgate) OnInit(app module.App, settings *conf.ModuleSettings) {
	self.BaseModule.OnInit(self, app, settings)
	self.SetListener(self)
}

func (self *httpgate) startHttpServer() *http.Server {
	srv := &http.Server{
		Addr:    ":8090",
		Handler: httpgateway.NewHandler(self.App),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			// cannot panic, because this probably is an intentional close
			log.Info("Httpserver: ListenAndServe() error: %s", err)
		}
	}()
	// returning reference so caller can call Shutdown()
	return srv
}

func (self *httpgate) Run(closeSig chan bool) {
	log.Info("httpgate: starting HTTP server :8090")
	srv := self.startHttpServer()
	<-closeSig
	log.Info("httpgate: stopping HTTP server")
	// now close the server gracefully ("shutdown")
	// timeout could be given instead of nil as a https://golang.org/pkg/context/
	if err := srv.Shutdown(nil); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}
	log.Info("httpgate: done. exiting")
}

func (self *httpgate) OnDestroy() {
	//别忘了继承
	self.BaseModule.OnDestroy()
}

//--------httpgateway
func (self *httpgate) httpgateway(request *go_api.Request) (*go_api.Response, error) {
	mux := http.NewServeMux()
	mux.HandleFunc("/httpgate/topic", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte(`hello world`))
	})

	req, err := http.NewRequest(request.Method, request.Url, strings.NewReader(request.Body))
	if err != nil {
		return nil, err
	}
	for _, v := range request.Header {
		req.Header.Set(v.Key, strings.Join(v.Values, ","))
	}
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	resp := &go_api.Response{
		StatusCode: int32(rr.Code),
		Body:       rr.Body.String(),
		Header:     make(map[string]*go_api.Pair),
	}
	for key, vals := range rr.Header() {
		header, ok := resp.Header[key]
		if !ok {
			header = &go_api.Pair{
				Key: key,
			}
			resp.Header[key] = header
		}
		header.Values = vals
	}
	return resp, nil
}
func (self *httpgate) NoFoundFunction(fn string) (*mqrpc.FunctionInfo, error) {
	return &mqrpc.FunctionInfo{
		Function:  reflect.ValueOf(self.httpgateway),
		Goroutine: true,
	}, nil
}
func (self *httpgate) BeforeHandle(fn string, callInfo *mqrpc.CallInfo) error {
	return nil
}
func (self *httpgate) OnTimeOut(fn string, Expired int64) {

}
func (self *httpgate) OnError(fn string, callInfo *mqrpc.CallInfo, err error) {}
func (self *httpgate) OnComplete(fn string, callInfo *mqrpc.CallInfo, result *rpcpb.ResultInfo, exec_time int64) {
}
