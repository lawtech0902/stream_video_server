package main

import (
	"github.com/julienschmidt/httprouter"
	"go_projects/stream_video_server/scheduler/task_runner"
	"net/http"
)

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()

	router.GET("/video-delete-record/:vid-id", vidDelRecHandler)

	return router
}

func main() {
	go task_runner.Start()
	r := RegisterHandlers()
	http.ListenAndServe(":9001", r)
}
