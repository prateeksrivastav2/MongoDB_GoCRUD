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

	// Connect to MongoDB
	Mongo.ConnectMongo()
	if Mongo.DB == nil {
		fmt.Println("Failed to connect to MongoDB.")
		return
	}

	// Initialize user collection
	dbase := Mongo.DB.Database("Users")
	userCollection := dbase.Collection("users")
	controllers.InitializeUserDatabase(userCollection)

	// Set up routes
	r := mux.NewRouter()
	routes.SetupRoutes(r)

	// Start the server
	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", r)
}