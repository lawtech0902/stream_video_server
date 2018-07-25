package main

import (
	"net/http"
	"go_projects/stream_video_server/api/defs"
	"encoding/json"
	"io"
)

func sendErrorResponse(w http.ResponseWriter, errResp defs.ErrorResponse) {
	w.WriteHeader(errResp.HttpSC)
	res, _ := json.Marshal(&errResp.Error)
	io.WriteString(w, string(res))
}

func sendNormalResponse(w http.ResponseWriter, resp string, sc int) {
	w.WriteHeader(sc)
	io.WriteString(w, resp)
}
