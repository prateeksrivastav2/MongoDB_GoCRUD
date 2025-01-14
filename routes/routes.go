package routes

import (
	controllers "g_o/controllers"

	"github.com/gorilla/mux"
)

// SetupRoutes configures the API routes.
func SetupRoutes(r *mux.Router) {
	r.HandleFunc("/user", controllers.CreateUserHandler).Methods("POST")
	r.HandleFunc("/users", controllers.CreateUsersHandler).Methods("POST")
	r.HandleFunc("/users", controllers.GetAllUserHandler).Methods("GET")
	r.HandleFunc("/user/{empid}", controllers.GetUserHandler).Methods("GET")
	r.HandleFunc("/user/{empid}", controllers.UpdateUserHandler).Methods("PUT")
	r.HandleFunc("/user/{empid}", controllers.DeleteUserHandler).Methods("DELETE")
}
