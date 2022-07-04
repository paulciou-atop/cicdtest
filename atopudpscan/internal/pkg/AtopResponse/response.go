package AtopResponse

import (
	"net/http"

	"gopkg.in/mgo.v2/bson"
)

type AtopResponse struct {
	Body   bson.M
	Writer http.ResponseWriter
}

func (r *AtopResponse) SendResponse(j []byte) {
	w := r.Writer
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(j)
}

func NewResponse(r *http.Request, w http.ResponseWriter) AtopResponse {
	var response = bson.M{}
	return AtopResponse{Body: response, Writer: w}
}
