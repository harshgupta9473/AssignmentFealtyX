package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	// "strings"

	// "log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/harshgupta9473/GoAssignmentFealtyX/middleware"
	"github.com/harshgupta9473/GoAssignmentFealtyX/models"
)

type StudentHandler struct {
}

func NewStudentHandler() *StudentHandler {
	return &StudentHandler{
	}
}

func (sh *StudentHandler) CreateStudent(w http.ResponseWriter, r *http.Request) {
	// sh.l.Println("creating a Student")
	student := r.Context().Value(middleware.UserContextKey).(models.Student)
	models.CreateStudent(&student)
}

func (sh *StudentHandler) GetAllStudent(w http.ResponseWriter, r *http.Request) {
	studetsList := models.GetStudentsList()
	w.Header().Set(`Content-Type`, `application/json`)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(studetsList)
}

func (sh *StudentHandler) GetAStudentByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idstring := vars["id"]

	id, err := strconv.Atoi(idstring)
	if err != nil {
		http.Error(w, "error converting idstring to int", http.StatusInternalServerError)
		return
	}
	student, _, err := models.GetStudentByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set(`Content-Type`, `application/json`)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(student)
}

func (sh *StudentHandler) UpdateAStudent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idstring := vars["id"]
	id, err := strconv.Atoi(idstring)
	if err != nil {
		http.Error(w, "error converting idstring to int", http.StatusInternalServerError)
		return
	}
	student := r.Context().Value(middleware.UserContextKey).(models.Student)
	err = models.UpdateStudent(id, &student)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (sh *StudentHandler) DeleteStudentByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idstring := vars["id"]
	id, err := strconv.Atoi(idstring)
	if err != nil {
		http.Error(w, "error converting idstring to int", http.StatusInternalServerError)
		return
	}
	err = models.DeleteStudentByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (sh *StudentHandler) GenerateStudentSummary(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idstring := vars["id"]
	id, err := strconv.Atoi(idstring)
	if err != nil {
		http.Error(w, "error converting idstring to int", http.StatusInternalServerError)
		return
	}
	student,_,err:=models.GetStudentByID(id)
	if err!=nil{
		http.Error(w,"student does not exists",http.StatusBadRequest)
		return
	}

	summary,err:=generateStudentSummaryOLLAMA(student)
	if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	// summary:=strings.ReplaceAll(summaryRes,"\\n",`\n`)
	w.Header().Set(`Content-Type`,`applicatio/json`)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"summary":summary,
	})
}

type OllamaReq struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool `json:"stream"`
}

func generateStudentSummaryOLLAMA(student *models.Student) (string,error){
	var summary string
	reqBody := OllamaReq{
		Model:  "llama3",
		Prompt: fmt.Sprintf("Generate a summary for the following student:Name:%s Age:%d Email:%s", student.Name, student.Age, student.Email),
		Stream: false,
	}
	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		return summary, fmt.Errorf("unable to marshal the request")
	}
	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(reqJSON))
	if err != nil {
		return summary, fmt.Errorf("error in connecting to LLM")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err!=nil{
		return "",err
	}
	var responseChunks map[string]interface{}
	err= json.Unmarshal(body, &responseChunks)
	if err != nil {
		return summary, fmt.Errorf("unable to unmarshal the response")
	}
	 
	if response,exists:=responseChunks["response"];exists{
		return response.(string),nil
	}else{
		return "",fmt.Errorf("unable to generate summary")
	}
}
