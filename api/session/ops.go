package session

import (
	"sync"
	"go_projects/stream_video_server/api/dbops"
	"go_projects/stream_video_server/api/defs"
	"go_projects/stream_video_server/api/utils"
	"time"
	"fmt"
)

var sessionMap *sync.Map

func nowInMilli() int64 {
	return time.Now().UnixNano() / 1000000
}

func deleteExpiredSession(sid string) {
	sessionMap.Delete(sid)
	dbops.DeleteSession(sid)
}

func init() {
	sessionMap = &sync.Map{}
}

func LoadSessionsFromDB() {
	r, err := dbops.RetrieveAllSessions()
	if err != nil {
		return
	}

	r.Range(func(k, v interface{}) bool {
		ss := v.(*defs.SimpleSession)
		sessionMap.Store(k, ss)
		return true
	})

	return
}

func GenerateNewSessionId(un string) string {
	id, _ := utils.NewUUID()
	ct := nowInMilli()
	ttl := ct + 30*60*1000 // 30 min
	ss := &defs.SimpleSession{
		UserName: un,
		TTL:      ttl,
	}
	sessionMap.Store(id, ss)

	err := dbops.InsertSession(id, ttl, un)
	if err != nil {
		return fmt.Sprintf("Error of GenerateNewSessionId: %s", err)
	}
	return id
}

func IsSessionExpired(sid string) (string, bool) {
	ss, ok := sessionMap.Load(sid)
	if ok {
		ct := nowInMilli()
		if ss.(*defs.SimpleSession).TTL < ct {
			deleteExpiredSession(sid)
			return "", true
		}

		return ss.(*defs.SimpleSession).UserName, false
	}

	return "", true
}
