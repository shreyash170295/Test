package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/mux"
	"github.com/practice_methods/student_project/student"
)

var students map[string]student.Student

func getStudents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var allStudents student.AllStudents
	allStudents.Students = []*student.Student{}
	for _, s := range students {
		allStudents.Students = append(allStudents.Students, &s)
	}
	b, err := proto.Marshal(&allStudents)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println(string(b))
	w.Write(b)
}

func getStudent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Write([]byte(`{"message":"incorrect Method expected get"}`))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "appication/json")
	recID := mux.Vars(r)["id"]
	if len(recID) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	s, ok := students[recID]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	b, err := json.Marshal(s)
	if b == nil {
		w.WriteHeader(http.StatusBadRequest)
	}
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
	var student student.Student
	err = json.Unmarshal(b, &student)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	id := strconv.Itoa(rand.Intn(2000))
	student.Id = id
	students[id] = student

	b, err = json.MarshalIndent(students[id], "", " ")
	if b == nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(b)
}

func updateStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	recID := mux.Vars(r)["id"]
	if len(recID) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	b, err := ioutil.ReadAll(r.Body)
	if b == nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var student student.Student
	err = json.Unmarshal(b, &student)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	students[recID] = student

	b, err = json.Marshal(students)
	if b == nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(b)
}

func deleteStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	recID := mux.Vars(r)["id"]
	if len(recID) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	delete(students, recID)

	b, err := json.Marshal(students)
	if b == nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	router := mux.NewRouter()

	students = make(map[string]student.Student)

	students["1"] = student.Student{Id: "1", Class: "10th", Section: "A", Name: &student.Name{Firstname: "Akash", Lastname: "Ghate"}}
	students["2"] = student.Student{Id: "2", Class: "10th", Section: "A", Name: &student.Name{Firstname: "Nik", Lastname: "Irving"}}
	students["3"] = student.Student{Id: "3", Class: "10th", Section: "B", Name: &student.Name{Firstname: "Josh", Lastname: "Gates"}}

	router.HandleFunc("/api/students", getStudents).Methods("GET").Name("GetAllStudents")
	router.HandleFunc("/api/students/{id}", getStudent).Methods("GET").Name("GetStudentByID")
	router.HandleFunc("/api/students", createStudent).Methods("POST").Name("AddStudentByID")
	router.HandleFunc("/api/students/{id}", updateStudent).Methods("PUT").Name("UpdateStudentByID")
	router.HandleFunc("/api/students/{id}", deleteStudent).Methods("DELETE").Name("DeleteStudentByID")

	log.Fatal(http.ListenAndServe(":8000", router))

}
