package main

import (
	"fmt"
	"g_o/Mongo"
	"g_o/controllers"

	// Mongo "g_o/mongo"

	//  "g_o/mongo"
	"g_o/routes"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("Crud in MongoDB")

	// connect to mongo
	Mongo.ConnectMongo()
	if Mongo.DB == nil {
		fmt.Println("Failed to connect to MongoDB.")
		return
	}

	dbase := Mongo.DB.Database("Users")
	userCollection := dbase.Collection("users")
	controllers.InitializeUserDatabase(userCollection)

	// routes
	r := mux.NewRouter()
	routes.SetupRoutes(r)

	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
