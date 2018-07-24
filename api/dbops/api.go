package dbops

import (
	_ "github.com/go-sql-driver/mysql"
	"log"
	"database/sql"
	"go_projects/stream_video_server/api/defs"
	"go_projects/stream_video_server/api/utils"
	"time"
)

// users
// insert user_info
func AddUserCredential(loginName, pwd string) error {
	stmtIns, err := dbConn.Prepare(`INSERT INTO users (login_name, pwd) VALUES (?, ?)`)
	if err != nil {
		return err
	}

	_, err = stmtIns.Exec(loginName, pwd)
	if err != nil {
		return err
	}

	defer stmtIns.Close()
	return nil
}

// query pwd
func GetUserCredential(loginName string) (string, error) {
	stmtOut, err := dbConn.Prepare(`SELECT pwd FROM users WHERE login_name = ?`)
	if err != nil {
		log.Printf("%s", err)
		return "", nil
	}

	var pwd string
	err = stmtOut.QueryRow(loginName).Scan(&pwd)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}

	defer stmtOut.Close()
	return pwd, nil
}

// delete user_info
func DeleteUser(loginName, pwd string) error {
	stmtDel, err := dbConn.Prepare(`DELETE FROM users WHERE login_name = ? AND pwd = ?`)
	if err != nil {
		log.Printf("Delete user error: %s", err)
		return err
	}

	_, err = stmtDel.Exec(loginName, pwd)
	if err != nil {
		return err
	}

	defer stmtDel.Close()
	return nil
}

// video_info
// insert video_info
func AddVideoInfo(aid int, name string) (*defs.VideoInfo, error) {
	// create uuid
	vid, err := utils.NewUUID()
	if err != nil {
		return nil, err
	}

	t := time.Now()
	ctime := t.Format("Jan 02 2006, 15:04:05") // must use this layout

	stmtIns, err := dbConn.Prepare(`INSERT INTO video_Info (video_id, author_id, name, display_ctime) VALUES (?, ?, ?, ?)`)
	if err != nil {
		return nil, err
	}
	_, err = stmtIns.Exec(vid, aid, name, ctime)
	if err != nil {
		return nil, err
	}

	res := &defs.VideoInfo{
		Id:           vid,
		AuthorId:     aid,
		Name:         name,
		DisplayCtime: ctime,
	}
	defer stmtIns.Close()
	return res, nil
}

// get video_info by vid
func GetVideoInfo(vid string) (*defs.VideoInfo, error) {
	stmtOut, err := dbConn.Prepare(`SELECT author_id, name, display_ctime FROM video_info WHERE video_id = ?`)
	if err != nil {
		return nil, err
	}

	var (
		aid  int
		name string
		dsc  string
	)
	err = stmtOut.QueryRow(vid).Scan(&aid, &name, &dsc)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	defer stmtOut.Close()
	res := &defs.VideoInfo{
		Id:           vid,
		AuthorId:     aid,
		Name:         name,
		DisplayCtime: dsc,
	}
	return res, nil
}

// delete video_info
func DeleteVideoInfo(vid string) error {
	stmtDel, err := dbConn.Prepare(`DELETE FROM video_info WHERE video_id = ?`)
	if err != nil {
		return err
	}
	_, err = stmtDel.Exec(vid)
	if err != nil {
		return err
	}
	defer stmtDel.Close()
	return nil
}

// add comment
func AddNewComments(vid string, aid int, content string) error {
	id, err := utils.NewUUID()
	if err != nil {
		return err
	}

	stmtIns, err := dbConn.Prepare(`INSERT INTO comments (id, video_id, author_id, content) VALUES (?, ?, ?, ?)`)
	if err != nil {
		return err
	}

	_, err = stmtIns.Exec(id, vid, aid, content)
	if err != nil {
		return err
	}

	defer stmtIns.Close()
	return nil
}

// list comment
func ListComments(vid string, from, to int) ([]*defs.Comment, error) {
	stmtOut, err := dbConn.Prepare(`SELECT comments.id, users.login_name, comments.content FROM comments
		INNER JOIN users ON comments.author_id = users.id
		WHERE comments.video_id = ? AND comments.time >  FROM_UNIXTIME(?) AND comments.time<=FROM_UNIXTIME(?)`)

	if err != nil {
		return nil, err
	}

	var res []*defs.Comment

	rows, err := stmtOut.Query(vid, from, to)
	if err != nil {
		return res, nil
	}

	for rows.Next() {
		var id, name, content string
		if err := rows.Scan(&id, &name, &content); err != nil {
			return res, err
		}

		c := &defs.Comment{
			Id:      id,
			VideoId: vid,
			Author:  name,
			Content: content,
		}
		res = append(res, c)
	}

	defer stmtOut.Close()
	return res, nil
}
