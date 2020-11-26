package main

import (
	"encoding/json"
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
	json.Marshal(students)

}

func getStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "appication/json")
	params := mux.Vars(r)

	//Loop through books and find id
	for _, item := range students {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Student{})

}

func createStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "appication/json")

	var student Student
	_ = json.NewDecoder(r.Body).Decode(&student)
	student.ID = strconv.Itoa(rand.Intn(1000))
	students[student.ID] = student
	json.NewEncoder(w).Encode(student)
}

func updateStudent(w http.ResponseWriter, r *http.Request) {

}

func deleteStudent(w http.ResponseWriter, r *http.Request) {

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
