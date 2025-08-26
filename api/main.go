package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	firebase "firebase.google.com/go/v4"
	"github.com/gorilla/mux"
	"google.golang.org/api/option"
	"cloud.google.com/go/firestore"
)

var client *firestore.Client

// Initialize Firebase
func initFirebase() {
	ctx := context.Background()
	opt := option.WithCredentialsFile("serviceAccountKey.json") // path to your service account key JSON

	// Use your Firebase Project ID (from the JS config you pasted)
	conf := &firebase.Config{ProjectID: "playitloud-1e8fe"}

	app, err := firebase.NewApp(ctx, conf, opt)
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

// POST /users → Add user
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

// GET /users → Fetch all users
func getUsers(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var users []User

	iter := client.Collection("users").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		var u User
		doc.DataTo(&u)
		users = append(users, u)
	}

	json.NewEncoder(w).Encode(users)
}

func main() {
	initFirebase()
	defer client.Close()

	r := mux.NewRouter()
	r.HandleFunc("/users", addUser).Methods("POST")
	r.HandleFunc("/users", getUsers).Methods("GET")

	fmt.Println("🚀 Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
