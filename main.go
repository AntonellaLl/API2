package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Users struct {
	Name     string    `json:"name"`
	LastName string    `json:"last_name"`
	CreateAt time.Time `json:"create_at"`
}

var usersStore = make(map[string]Users)
var id int

//GetUsers - GET - /api/notes
func GetUsers(w http.ResponseWriter, r *http.Request) {
	var usuarios []Users
	for _, v := range usersStore {
		usuarios = append(usuarios, v)
	}
	w.Header().Set("Content-Type", "aplication/jason")
	j, err := json.Marshal(usuarios)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

//POSTUsers - POST - /api/notes

func POSTUsers(w http.ResponseWriter, r *http.Request) {
	var usuarios Users
	err := json.NewDecoder(r.Body).Decode(&usuarios)
	if err != nil {
		panic(err)
	}
	usuarios.CreateAt = time.Now()
	id++
	k := strconv.Itoa(id)
	usersStore[k] = usuarios

	w.Header().Set("Content-Type", "aplication/jason")
	j, err := json.Marshal(usuarios)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}

//PUTUsers - PUT - /api/notes

func PUTUsers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	k := vars["id"]
	var usersUpdate Users
	err := json.NewDecoder(r.Body).Decode(&usersUpdate)
	if err != nil {
		panic(err)
	}
	if usuarios, ok := usersStore[k]; ok {
		usersUpdate.CreateAt = usuarios.CreateAt
		delete(usersStore, k)
		usersStore[k] = usersUpdate
	} else {
		log.Printf(" no encontramos el id %s", k)
	}
	w.WriteHeader(http.StatusNoContent)
}

//DELETEUsers - DELETE - /api/users

func DELETEUsers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	k := vars["id"]
	if _, ok := usersStore[k]; ok {
		delete(usersStore, k)
	} else {
		log.Printf(" no encontramos el id %s", k)
	}
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	r := mux.NewRouter().StrictSlash(false)
	r.HandleFunc("/api/users", GetUsers).Methods("GET")
	r.HandleFunc("/api/users", POSTUsers).Methods("POST")
	r.HandleFunc("/api/users/{id}", PUTUsers).Methods("PUT")
	r.HandleFunc("/api/users/{id}", DELETEUsers).Methods("DELETE")

	server := &http.Server{
		Addr:           ":8080",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Println("Listening http://localhost:8080...")
	server.ListenAndServe()

}
