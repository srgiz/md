package internal_test // fix: circle import

import (
	"md/internal/app/server"
	"testing"
)

func TestRequestHandler(t *testing.T) {
	s := server.Initialize("/tmp/data/")
	w := s.TestJsonRpc("File.Test", nil)

	assertJsonRpcErrorCode(t, w, -32000)
	assertJsonRpcId(t, w, "test")
	assertResponseContains(t, w, `"rpc: can't find method \"File.Test\""`)
}
