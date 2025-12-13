package http

import (
	"encoding/json"
	"log"
	"md/internal/io/http/internal"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gorilla/rpc/v2"
	"github.com/gorilla/rpc/v2/json2"
)

type server struct {
	rpcServer *rpc.Server
}

func newServer(
	fileReceiver *internal.FileReceiver,
) *server {
	rpcServer := rpc.NewServer()
	rpcServer.RegisterCodec(json2.NewCustomCodec(&rpc.CompressionSelector{}), "application/json")

	addHandler(rpcServer, fileReceiver, "File")

	return &server{rpcServer: rpcServer}
}

func addHandler(rpcServer *rpc.Server, receiver any, name string) {
	if err := rpcServer.RegisterService(receiver, name); err != nil {
		log.Fatal(err)
	}
}

func (app *server) Run() {
	//http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./"))))
	http.Handle("/jsonrpc", app.rpcServer)

	httpServer := &http.Server{
		Addr: ":8080",
		//ReadHeaderTimeout: 5 * time.Minute,
	}

	log.Fatal(httpServer.ListenAndServe())
}

func (app *server) Test(r *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	app.rpcServer.ServeHTTP(w, r)
	return w
}

func (app *server) TestJsonRpc(method string, params any) *httptest.ResponseRecorder {
	payload := map[string]any{"jsonrpc": "2.0", "method": method, "params": params, "id": "test"}
	body, _ := json.Marshal(payload)
	return app.Test(httptest.NewRequest(http.MethodPost, "/jsonrpc", strings.NewReader(string(body))))
}
