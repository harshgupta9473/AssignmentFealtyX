package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/harshgupta9473/GoAssignmentFealtyX/models"
)
type ContextKey string

const UserContextKey ContextKey = "student"

func MiddlewareValidation(next http.Handler)http.Handler{
	return http.HandlerFunc(func (w http.ResponseWriter,r *http.Request)  {
		var student models.Student
		err:=json.NewDecoder(r.Body).Decode(&student)
		if err!=nil{
			http.Error(w,"not correct formate",http.StatusBadRequest)
			return
		}
		err=student.Validate()
		if err!=nil{
			http.Error(w,"error validating the student detail, check again before submitting",http.StatusBadRequest)
			return
		}
		ctx:=context.WithValue(r.Context(),UserContextKey,student)
		next.ServeHTTP(w,r.WithContext(ctx))
	})
}