package main

import (
	"net/http"
	"go_projects/stream_video_server/api/session"
	"go_projects/stream_video_server/api/defs"
)

var (
	HEADER_FIELD_SESSION = "X-session-Id"
	HEADER_FIELD_UNAME   = "X-User-Name"
)

func validateUserSession(r *http.Request) bool {
	sid := r.Header.Get(HEADER_FIELD_SESSION)
	if len(sid) == 0 {
		return false
	}

	uname, ok := session.IsSessionExpired(sid)
	if ok {
		return false
	}

	r.Header.Add(HEADER_FIELD_UNAME, uname)
	return true
}

func validateUser(w http.ResponseWriter, r *http.Request) bool {
	uname := r.Header.Get(HEADER_FIELD_UNAME)
	if len(uname) == 0 {
		sendErrorResponse(w, defs.ErrorNotAuthUser)
		return false
	}

	return true
}
