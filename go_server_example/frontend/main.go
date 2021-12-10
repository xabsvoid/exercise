package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

type PagePatch struct {
	Backend string
	ID      uint64
	Limit   uint16
}

var backendUrl string

var templates = template.Must(template.ParseFiles(
	"tmpl/delete.html",
	"tmpl/get.html",
	"tmpl/post.html",
	"tmpl/put.html"))

func getUint64(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/getNotes/", http.StatusMovedPermanently)
}

func getNotesHandler(w http.ResponseWriter, r *http.Request) {
	var patch PagePatch
	patch.Backend = backendUrl
	patch.ID = 0
	patch.Limit = 20
	err := templates.ExecuteTemplate(w, "get.html", patch)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func postNoteHandler(w http.ResponseWriter, r *http.Request) {
	var patch PagePatch
	patch.Backend = backendUrl
	err := templates.ExecuteTemplate(w, "post.html", patch)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func putNoteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := getUint64(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var patch PagePatch
	patch.Backend = backendUrl
	patch.ID = id
	err = templates.ExecuteTemplate(w, "put.html", patch)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func deleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := getUint64(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var patch PagePatch
	patch.Backend = backendUrl
	patch.ID = id
	err = templates.ExecuteTemplate(w, "delete.html", patch)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	addr := os.Getenv("FRONTEND_ADDR")
	backendUrl = os.Getenv("FRONTEND_BACKEND_URL")
	router := mux.NewRouter()
	router.HandleFunc("/", rootHandler)
	router.HandleFunc("/getNotes/", getNotesHandler)
	router.HandleFunc("/postNote/", postNoteHandler)
	router.HandleFunc("/putNote/{id:[0-9]+}", putNoteHandler)
	router.HandleFunc("/deleteNote/{id:[0-9]+}", deleteNoteHandler)
	log.Fatal(http.ListenAndServe(addr, router))
}
