package internal_test // fix: circle import

import (
	"fmt"
	"md/internal/app/server"
	"testing"
)

func TestUnknownMethod(t *testing.T) {
	s := server.Initialize("/tmp/data/")
	w := s.TestJsonRpc("File.Test", nil)

	assertJsonRpcErrorCode(t, w, -32000)
	assertJsonRpcId(t, w, "test")
	assertResponseContains(t, w, `"rpc: can't find method \"File.Test\""`)
}

func TestValidFilepath(t *testing.T) {
	s := server.Initialize("/tmp/data/")

	var tests = []struct {
		path   string
		substr string
	}{
		{
			path:   "/foo",
			substr: `{"jsonrpc":"2.0","result":{"id":`,
		},
		{
			path:   ".bar",
			substr: "Field validation for 'Path' failed on the 'allowedFilepath' tag",
		},
		{
			path:   "/.baz",
			substr: "Field validation for 'Path' failed on the 'allowedFilepath' tag",
		},
		{
			path:   ".folder/foo",
			substr: "Field validation for 'Path' failed on the 'allowedFilepath' tag",
		},
		{
			path:   "/foo/baz",
			substr: `{"jsonrpc":"2.0","result":{"id":`,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Path: %s", test.path), func(t *testing.T) {
			w := s.TestJsonRpc("File.Find", map[string]any{"path": test.path})

			assertResponseContains(t, w, test.substr)
		})
	}
}
