package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"io"
	"io/ioutil"
	"go_projects/stream_video_server/api/defs"
	"encoding/json"
	"go_projects/stream_video_server/api/dbops"
	"go_projects/stream_video_server/api/session"
)

func CreateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// read body from request
	res, _ := ioutil.ReadAll(r.Body)
	uBody := &defs.UserCredential{}

	if err := json.Unmarshal(res, uBody); err != nil {
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
		return
	}

	if err := dbops.AddUserCredential(uBody.UserName, uBody.Pwd); err != nil {
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}

	id := session.GenerateNewSessionId(uBody.UserName)
	su := &defs.SignedUp{
		Success:   true,
		SessionId: id,
	}

	if resp, err := json.Marshal(su); err != nil {
		sendErrorResponse(w, defs.ErrorInternalFaults)
		return
	} else {
		sendNormalResponse(w, string(resp), 200)
	}
}

func Login(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	uname := p.ByName("user_name")
	io.WriteString(w, uname)
}
