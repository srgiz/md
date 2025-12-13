package internal_test

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"strings"

	"github.com/gorilla/rpc/v2/json2"
	"github.com/stretchr/testify/assert"
)

var responses = make(map[*httptest.ResponseRecorder]struct {
	body []byte
})

func parseResponse(w *httptest.ResponseRecorder) {
	if _, ok := responses[w]; !ok {
		result := w.Result()
		defer result.Body.Close()

		bb, _ := io.ReadAll(result.Body)
		responses[w] = struct{ body []byte }{body: bb}
	}
}

func assertResponseContains(t assert.TestingT, response *httptest.ResponseRecorder, substr string) {
	parseResponse(response)
	assert.True(t, strings.Contains(string(responses[response].body), substr))
}

func jsonRpcErrorResponse(w *httptest.ResponseRecorder) *struct {
	Error *json2.Error
	Id    string
} {
	parseResponse(w)

	v := &struct {
		Error *json2.Error
		Id    string
	}{}

	json.Unmarshal(responses[w].body, v)
	return v
}

func assertJsonRpcErrorCode(t assert.TestingT, response *httptest.ResponseRecorder, code int) {
	assert.Equal(t, int(jsonRpcErrorResponse(response).Error.Code), code)
}

func assertJsonRpcId(t assert.TestingT, response *httptest.ResponseRecorder, id string) {
	expectedId := jsonRpcErrorResponse(response).Id
	assert.Greater(t, len(expectedId), 0)
	assert.Equal(t, expectedId, id)
}
