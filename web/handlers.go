package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"log"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/url"
	"net/http/httputil"
)

type HomePage struct {
	Name string
}

type UserPage struct {
	Name string
}

func homeHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	cname, err1 := r.Cookie("username")
	sid, err2 := r.Cookie("session")
	if err1 != nil || err2 != nil {
		pg := &HomePage{Name: "lawtech"}
		t, err := template.ParseFiles("./templates/home.html")
		if err != nil {
			log.Printf("Parsing template home.html error: %v", err)
			return
		}

		t.Execute(w, pg)
		return
	}

	if len(cname.Value) != 0 && len(sid.Value) == 0 {
		http.Redirect(w, r, "/userhome", http.StatusFound)
		return
	}
}

func userHomeHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	cname, err1 := r.Cookie("username")
	_, err2 := r.Cookie("session")
	if err1 != nil || err2 != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	fname := r.FormValue("username")
	var pg *UserPage
	if len(cname.Value) != 0 {
		pg = &UserPage{Name: cname.Value} // load from session if authenticated
	} else if len(fname) != 0 {
		pg = &UserPage{Name: fname}
	}

	t, err := template.ParseFiles("./templates/userhome.html")
	if err != nil {
		log.Printf("Parsing template userhome.html error: %v", err)
		return
	}
	t.Execute(w, pg)
}

func apiHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// pretreatment before deal with request
	if r.Method != http.MethodPost {
		res, _ := json.Marshal(ErrorRequestNotRecognized)
		io.WriteString(w, string(res))
		return
	}

	res, _ := ioutil.ReadAll(r.Body)
	apiBody := &ApiBody{}
	if err := json.Unmarshal(res, apiBody); err != nil {
		res, _ := json.Marshal(ErrorRequestBodyParseFailed)
		io.WriteString(w, string(res))
		return
	}

	request(apiBody, w, r)
	defer r.Body.Close()
}

func proxyHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	u, _ := url.Parse("http://127.0.0.1:9000/")
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ServeHTTP(w, r)
}
