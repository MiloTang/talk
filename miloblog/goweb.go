package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"miloblog/gocs"
	"miloblog/wr"
	"net/http"
	"os"
	"strings"
)

type Page struct {
	Title    string
	Lists    []Context
	Next     string
	Previous string
	Token    string
	Info     string
	Details  string
	Username string
	Password string
	Error    error
}
type Context struct {
	Introduction string
	Link         string
}

var (
	p                         = &Page{}
	cs    *gocs.CookieSession = nil
	Debug bool                = true
	err   error               = nil
	file  *os.File
)

func entry(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		wr.Index(w, r)
	} else {
		urls := strings.Split(r.URL.Path, "/")
		switch urls[1] {
		case "life":
			wr.Life(w, r)
		case "admin":
			wr.AdminLogin(w, r, cs)
		case "manuallist":
			wr.Manual(w, r)
		case "index":
			wr.Index(w, r)
		case "blog", "manual":
			wr.Details(w, r, urls[1])
		case "editor":
			wr.Editor(w, r, cs)
		case "delsession":
			wr.Delsession(w, r, cs)
		case "upload":
			upload(w, r)
		default:
			ErrorPage(w, r)
		}
	}
}
func ErrorPage(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("error.html")
	t.Execute(w, nil)
}
func main() {
	gocs.MaxLT = 3600
	gocs.CookieName = "Miloblog"
	cs = gocs.NewCookieSession()
	http.Handle("/css/", http.FileServer(http.Dir("static")))
	http.Handle("/js/", http.FileServer(http.Dir("static")))
	http.Handle("/images/", http.FileServer(http.Dir("static")))
	http.Handle("/bootstrap/", http.FileServer(http.Dir("static")))
	http.Handle("/layoutitlib/", http.FileServer(http.Dir("static")))
	http.Handle("/wysiwyg/", http.FileServer(http.Dir("static")))
	http.Handle("/fonts/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/", entry)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
func upload(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Form)
	if r.Method == "GET" {
		t, err := template.ParseFiles("upload.html")
		checkErr(err)
		t.Execute(w, nil)
	} else {
		file, handle, err := r.FormFile("editorImage")
		checkErr(err)
		fileName := handle.Filename[strings.LastIndex(handle.Filename, "\\")+1:]
		f, err := os.OpenFile("./static/images/"+fileName, os.O_WRONLY|os.O_CREATE, 0666)
		io.Copy(f, file)
		checkErr(err)
		defer f.Close()
		defer file.Close()
		fmt.Println("upload success")
		fmt.Println("127.0.0.1/images/" + fileName)
		fmt.Fprintf(w, "http://127.0.0.1/images/"+fileName)
		return
	}
}
func checkErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}
