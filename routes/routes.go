package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/harshgupta9473/GoAssignmentFealtyX/handlers"
	"github.com/harshgupta9473/GoAssignmentFealtyX/middleware"
)

func RegisterRoutes(router *mux.Router, sh *handlers.StudentHandler) {

	router.Handle("/students",middleware.MiddlewareValidation(http.HandlerFunc(sh.CreateStudent))).Methods(http.MethodPost)
	router.Handle("/students/{id}",middleware.MiddlewareValidation(http.HandlerFunc(sh.UpdateAStudent))).Methods(http.MethodPut)
	
	router.HandleFunc("/students/{id}",sh.DeleteStudentByID).Methods(http.MethodDelete)
	router.HandleFunc("/students/{id}",sh.GetAStudentByID).Methods(http.MethodGet)
	router.HandleFunc("/students",sh.GetAllStudent).Methods(http.MethodGet)
	router.HandleFunc("/students/{id}/summary",sh.GenerateStudentSummary).Methods(http.MethodGet)
}
