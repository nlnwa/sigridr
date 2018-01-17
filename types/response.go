package types

import (
	"net/http"
)

type Response struct {
	Status   string              `json:"status"`
	Code     int                 `json:"code"`
	Protocol string              `json:"protocol"`
	Header   map[string][]string `json:"header"`
}

func (r *Response) FromHttpResponse(response *http.Response) *Response {
	r.Header = response.Header
	r.Status = response.Status
	r.Code = response.StatusCode
	r.Protocol = response.Proto
	return r
}
