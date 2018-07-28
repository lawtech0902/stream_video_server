#! /bin/bash

# build web ui
cd $GOPATH/src/go_projects/stream_video_server/web
go install
cp $GOPATH/bin/web $GOPATH/bin/video_server_web_ui/web
cp -r $GOPATH/src/go_projects/stream_video_server/templates $GOPATH/bin/video_server_web_ui/

