package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

//Student Model
type Student struct {
	ID      string `json: "id`
	Class   string `json: "class"`
	Section string `json: "section"`
	Name    *Name  `json: "name"`
}

//Name model
type Name struct {
	Firstname string `json: "firstname"`
	Lastname  string `json: "lastname"`
}

var students map[string]Student

func getStudents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.MarshalIndent(students, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func getStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "appication/json")
	recID := mux.Vars(r)["id"]
	s, ok := students[recID]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	b, err := json.Marshal(s)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func createStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "appication/json")
	b, err := ioutil.ReadAll(r.Body)
	if b == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var student Student
	err = json.Unmarshal(b, &student)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	id := strconv.Itoa(rand.Intn(2000))

	students[id] = student

	b, err1 := json.MarshalIndent(students, "", " ")
	if err1 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(b)
}

func updateStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	recID := mux.Vars(r)["id"]
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var student Student
	error := json.Unmarshal(b, &student)
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	students[recID] = student

	b, err = json.Marshal(students)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(b)
}

func deleteStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	recID := mux.Vars(r)["id"]

	delete(students, recID)

	b, err := json.Marshal(students)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func main() {
	router := mux.NewRouter()

	students = make(map[string]Student)

	students["1"] = Student{ID: "1", Class: "10th", Section: "A", Name: &Name{Firstname: "Akash", Lastname: "Ghate"}}
	students["2"] = Student{ID: "2", Class: "10th", Section: "A", Name: &Name{Firstname: "Nik", Lastname: "Irving"}}
	students["3"] = Student{ID: "3", Class: "10th", Section: "B", Name: &Name{Firstname: "Josh", Lastname: "Gates"}}

	router.HandleFunc("/api/students", getStudents).Methods("GET")
	router.HandleFunc("/api/student/{id}", getStudent).Methods("GET")
	router.HandleFunc("/api/student", createStudent).Methods("POST")
	router.HandleFunc("/api/student/{id}", updateStudent).Methods("PUT")
	router.HandleFunc("/api/student/{id}", deleteStudent).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))

}
