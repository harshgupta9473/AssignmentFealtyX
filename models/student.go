package models

import (
	"fmt"
	"sync"

	"github.com/go-playground/validator"
)

type Student struct {
	ID    int    `json:"id"`
	Name  string `json:"name" validate:"required"`
	Age   int    `json:"age" validate:"gte=4,lte=40"`
	Email string `json:"email" validate:"required,email"`
}

var (
	studentList = []*Student{}
	mu          sync.RWMutex
)

func GetStudentsList() []*Student {
	mu.RLock()
	defer mu.RUnlock()
	return studentList
}

func (st *Student) Validate() error {
	validate := validator.New()
	return validate.Struct(st)
}

func UpdateStudent(id int, student *Student) error {

	mu.Lock()
	defer mu.Unlock()
	_, index, err := getStudentByid(id)
	if err != nil {
		return err
	}
	student.ID = id
	studentList[index] = student
	return nil
}

func CreateStudent(student *Student) {
	mu.Lock()
	defer mu.Unlock()
	student.ID = getNextID()
	studentList = append(studentList, student)
}

func getNextID() int {
	var lp int
	if len(studentList) > 0 {
		lp = studentList[len(studentList)-1].ID
	} else {
		lp = 0
	}

	return lp + 1
}

func DeleteStudentByID(id int) error {
	mu.Lock()
	defer mu.Unlock()
	_, index, err := getStudentByid(id)
	if err != nil {
		return err
	}
	studentList = append(studentList[:index], studentList[index+1:]...)
	return nil
}

func getStudentByid(id int) (*Student, int, error) {

	for i, p := range studentList {
		if p.ID == id {
			return studentList[i], i, nil
		}
	}
	return nil, -1, fmt.Errorf("Student not found")
}

func GetStudentByID(id int) (*Student, int, error) {
	mu.RLock()
	defer mu.RUnlock()
	for i, p := range studentList {
		if p.ID == id {
			return studentList[i], i, nil
		}
	}
	return nil, -1, fmt.Errorf("Student not found")
}
