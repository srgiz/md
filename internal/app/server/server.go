package server

import (
	//"flag"
	"log"
	"md/internal/app/server/internal"
	"net/http"
	//"time"

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

	// File.Edit
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
