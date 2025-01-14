package controllers

import (
	"context"
	"encoding/json"
	"fmt"

	"g_o/Mongo"

	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection

func InitializeUserDatabase(Collection *mongo.Collection) {
	userCollection = Collection
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user Mongo.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created Successfully"})
}
func CreateUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var users []Mongo.User
	if err := json.NewDecoder(r.Body).Decode(&users); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var interfaceSlice []interface{}
	for _, user := range users {
		interfaceSlice = append(interfaceSlice, user)
	}
	// fmt.Println(interfaceSlice...)

	result, err := userCollection.InsertMany(ctx, interfaceSlice)
	if err != nil {
		http.Error(w, "Failed to insert users: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	response := map[string]interface{}{
		"message":      "Users created successfully",
		"inserted_ids": result.InsertedIDs,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to write response: "+err.Error(), http.StatusInternalServerError)
	}
}

func GetAllUserHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := userCollection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch users: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var users []Mongo.User
	if err := cursor.All(ctx, &users); err != nil {
		http.Error(w, "Failed to decode users: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Failed to write response: "+err.Error(), http.StatusInternalServerError)
	}
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	empid := vars["empid"]
	// fmt.Println(empid)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{Key: "emp_id", Value: empid}}

	var user Mongo.User
	err := userCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		http.Error(w, "Failed to fetch user: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(user)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Failed to write response: "+err.Error(), http.StatusInternalServerError)
	}
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	empid := vars["empid"]
	fmt.Println("yha tk")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{Key: "emp_id", Value: empid}}

	var updateData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	update := bson.D{{Key: "$set", Value: updateData}}

	updateResult, err := userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		http.Error(w, "Failed to update the user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"matchedCount":  updateResult.MatchedCount,
		"modifiedCount": updateResult.ModifiedCount,
		"message":       "User updated successfully",
	})
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	empid := vars["empid"]

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{Key: "emp_id", Value: empid}}

	result, err := userCollection.DeleteOne(ctx, filter)
	if err != nil {
		http.Error(w, "Failed to delete user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		http.Error(w, "No user found with the provided emp_id", http.StatusNotFound)
		return
	}

	response := map[string]string{"message": "User deleted successfully"}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to write response: "+err.Error(), http.StatusInternalServerError)
	}
}
