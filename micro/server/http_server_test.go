package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.houston.softwaregrp.net/onestack/micro/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mock struct {
}

func (*mock) GetAll() ([]db.Pet, error) {
	var pets []db.Pet
	var pet db.Pet
	pet.ID, _ = primitive.ObjectIDFromHex("61993177341fd4e12e2c734b")
	pet.Type = "Dog"
	pet.Breed = "Beagle"
	t, _ := time.Parse("2006-01-02", "2011-01-19")
	pet.BirthDate = &t
	return append(pets, pet), nil
}

func (*mock) Get(id primitive.ObjectID) (*db.Pet, error) {
	return nil, nil
}

func (*mock) Insert(*db.Pet) (primitive.ObjectID, error) {
	return primitive.ObjectID{}, nil
}

func (*mock) Delete(id primitive.ObjectID) (int64, error) {
	return 0, nil
}

func (*mock) Update(*db.Pet) (*db.Pet, error) {
	return nil, nil
}

func init() {
	iDB = &mock{}
}

func TestGetPet(t *testing.T) {
	req, err := http.NewRequest("GET", "/pet", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetPetAll)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `[{"_id":"61993177341fd4e12e2c734b","type":"Dog","breed":"Beagle","birthdate":"2011-01-19T00:00:00Z"}]`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
