package config

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

var FirebaseAuth *auth.Client

func InitFirebase() {
	opt := option.WithCredentialsFile("firebase-key.json") 
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Error initializing Firebase App: %v\n", err)
	}

	authClient, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("Error initializing Firebase Auth: %v\n", err)
	}

	FirebaseAuth = authClient
	log.Println("Firebase Auth initialized")
}
