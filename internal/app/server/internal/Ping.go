package internal

import "net/http"

type PingRequest struct {
}

type PingHandler struct {
}

func NewPingHandler() *PingHandler {
	return &PingHandler{}
}

// Ping.Test
func (h *PingHandler) Test(r *http.Request, req *PingRequest, res *PingResponse) error {
	*res = PingResponse{}
	res.Text = "pong"

	return nil
}

type PingResponse struct {
	Text string `json:"text"`
}
