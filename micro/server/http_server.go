package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.houston.softwaregrp.net/onestack/micro/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewMux() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/pet", GetPetAll).Methods("GET")
	r.HandleFunc("/pet/{id}", GetPet).Methods("GET")
	r.HandleFunc("/pet", CreatePet).Methods("POST")
	r.HandleFunc("/pet/{id}", DeletePet).Methods("DELETE")
	r.HandleFunc("/pet/{id}", EditPet).Methods("PUT")
	return r
}

func StartHttpServer() {
	srv := &http.Server{
		Handler:      NewMux(),
		Addr:         ":8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("Starting server")
	log.Fatal(srv.ListenAndServe())
}

func CreatePet(response http.ResponseWriter, request *http.Request) {
	b, err := io.ReadAll(request.Body)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(err.Error()))
		return
	}
	pet := db.Pet{}
	err = json.Unmarshal(b, &pet)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	}
	fmt.Println("inserting in db ", pet)
	res, err := pet.Insert()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	}
	response.Header().Add("Content-Type", "application/json")
	pet.ID = res
	json.NewEncoder(response).Encode(pet)
	log.Println(res)
}

func GetPetAll(response http.ResponseWriter, request *http.Request) {
	pet := db.Pet{}
	res, err := pet.GetAll()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
	}

	response.Header().Add("Content-Type", "application/json")
	json.NewEncoder(response).Encode(res)
	log.Println(res)
}

func GetPet(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	pet := db.Pet{}
	res, err := pet.Get(id)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
	}

	response.Header().Add("Content-Type", "application/json")
	json.NewEncoder(response).Encode(*res)
	log.Println(*res)
}

func DeletePet(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	pet := db.Pet{}
	res, err := pet.Delete(id)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
	}
	response.Header().Add("Content-Type", "application/json")
	str := fmt.Sprintf("resource %s deleted successfully", params["id"])
	json.NewEncoder(response).Encode(str)
	log.Println(res)
}

func EditPet(response http.ResponseWriter, request *http.Request) {
	b, err := io.ReadAll(request.Body)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(err.Error()))
	}
	pet := db.Pet{}
	err = json.Unmarshal(b, &pet)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	}
	vars := mux.Vars(request)
	pet.ID, _ = primitive.ObjectIDFromHex(vars["id"])
	res, err := pet.Update(pet.ID)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(response).Encode(*res)
	response.Header().Add("Content-Type", "application/json")
	log.Println(*res)
}
