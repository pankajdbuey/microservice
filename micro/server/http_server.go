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

type IDBOperation interface {
	Get(id primitive.ObjectID) (*db.Pet, error)
	GetAll() ([]db.Pet, error)
	Insert(db.Pet) (primitive.ObjectID, error)
	Delete(id primitive.ObjectID) (int64, error)
	Update(db.Pet) (int64, error)
}

var iDB IDBOperation

func init() {
	iDB = &db.Pet{}
}

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
	res, err := iDB.Insert(pet)
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
	res, err := iDB.GetAll()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	}

	response.Header().Add("Content-Type", "application/json")
	if len(res) == 0 {
		json.NewEncoder(response).Encode("No entry found")
		log.Println("No entry found")
	} else {
		json.NewEncoder(response).Encode(res)
		log.Println(res)
	}

}

func GetPet(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	res, err := iDB.Get(id)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	}

	response.Header().Add("Content-Type", "application/json")
	json.NewEncoder(response).Encode(*res)
	log.Println(*res)
}

func DeletePet(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	_, err := iDB.Delete(id)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	}
	response.Header().Add("Content-Type", "application/json")
	str := fmt.Sprintf("record %s deleted successfully", params["id"])
	json.NewEncoder(response).Encode(str)
	log.Println(str)
}

func EditPet(response http.ResponseWriter, request *http.Request) {
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
	vars := mux.Vars(request)
	pet.ID, _ = primitive.ObjectIDFromHex(vars["id"])
	_, err = iDB.Update(pet)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	}
	s := fmt.Sprintf("record %v modified successfully", vars["id"])
	json.NewEncoder(response).Encode(s)
	response.Header().Add("Content-Type", "application/json")
	log.Println(s)
}
