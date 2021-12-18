package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Server struct {
	db Database
}

func (server *Server) Run(dsn string, addr string) error {
	err := server.db.Init(dsn)
	if err != nil {
		return err
	}
	router := mux.NewRouter()
	router.Use(commonMiddleware)
	router.HandleFunc("/notes/{offset:[0-9]+}/{limit:[0-9]+}", server.getNotes).Methods("GET")
	router.HandleFunc("/note", server.postNote).Methods("POST")
	router.HandleFunc("/note/{id:[0-9]+}", server.getNote).Methods("GET")
	router.HandleFunc("/note/{id:[0-9]+}", server.updateNote).Methods("PUT")
	router.HandleFunc("/note/{id:[0-9]+}", server.deleteNote).Methods("DELETE")

	headersOk := handlers.AllowedHeaders([]string{"Access-Control-Allow-Origin", "Accept", "Accept-Language", "Content-Type", "Content-Language", "Origin"})
	originsOk := handlers.AllowedOrigins([]string{"http://localhost:*", "*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

	return http.ListenAndServe(addr, handlers.CORS(headersOk, originsOk, methodsOk)(router))
}

func commonMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		h.ServeHTTP(w, r)
	})
}

func getUint64(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}

func getUint16(s string) (uint16, error) {
	i, err := strconv.ParseUint(s, 10, 16)
	return uint16(i), err
}

func (server *Server) getNotes(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	offset, err := getUint64(vars["offset"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	limit, err := getUint16(vars["limit"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	notes, err := server.db.SelectNotes(offset, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(notes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (server Server) getNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := getUint64(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	note, err := server.db.SelectNote(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(note)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (server Server) postNote(w http.ResponseWriter, r *http.Request) {
	var note Note
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	note.CreatedTime = time.Now()
	err = note.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = server.db.InsertNote(note)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (server Server) updateNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := getUint64(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var note Note
	err = json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	note.ID = id
	err = note.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = server.db.UpdateNote(note)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (server Server) deleteNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := getUint64(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = server.db.DeleteNote(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
