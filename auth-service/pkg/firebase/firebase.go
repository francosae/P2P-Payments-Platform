package firebase

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

type FirebaseClient struct {
	App *firebase.App
}

func InitializeAppWithServiceAccount() (*FirebaseClient, error) {
	opt := option.WithCredentialsFile("./pkg/firebase/serviceAccountKey.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Error initializing Firebase Admin: %v\n", err)
	}
	return &FirebaseClient{App: app}, nil
}
