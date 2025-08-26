package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	firebase "firebase.google.com/go/v4"
	"cloud.google.com/go/firestore"
	"github.com/gorilla/mux"
	"google.golang.org/api/option"
)

var client *firestore.Client

// Initialize Firebase
func initFirebase() {
	ctx := context.Background()
	opt := option.WithCredentialsFile("serviceAccountKey.json") // Path to your Firebase service account key

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("Error initializing Firebase app: %v", err)
	}

	firestoreClient, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("Error initializing Firestore: %v", err)
	}

	client = firestoreClient
	fmt.Println("✅ Connected to Firestore!")
}

// Struct for incoming data
type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

// API handler to insert data
func addUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var user User

	// Decode JSON from request
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Insert into Firestore
	_, _, err = client.Collection("users").Add(ctx, map[string]interface{}{
		"name":  user.Name,
		"email": user.Email,
		"age":   user.Age,
	})
	if err != nil {
		http.Error(w, "Failed to add user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User added successfully"})
}

func main() {
	initFirebase()
	defer client.Close()

	r := mux.NewRouter()
	r.HandleFunc("/users", addUser).Methods("POST")

	fmt.Println("🚀 Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
