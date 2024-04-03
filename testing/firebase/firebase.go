package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

type firebaseAuthResponse struct {
	IDToken string `json:"idToken"`
}

type FirebaseClient struct {
	App *firebase.App
}

func InitializeAppWithServiceAccount() (*FirebaseClient, error) {
	opt := option.WithCredentialsFile("/Users/pfranco/Desktop/poolparty/services-workspace/testing/firebase/serviceAccountKey.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Error initializing Firebase Admin: %v\n", err)
	}
	return &FirebaseClient{App: app}, nil
}

func main() {
	var email, password string

	// fmt.Println("Enter Email:")
	// fmt.Scanln(&email)

	// fmt.Println("Enter Password:")
	// fmt.Scanln(&password)

	email = "paul@gmail.com"
	password = "123456"

	token, err := signInWithEmailPassword(email, password)
	if err != nil {
		log.Fatalf("Error signing in: %v\n", err)
	}

	client, err := InitializeAppWithServiceAccount()
	if err != nil {
		log.Fatalf("Error initializing Firebase Admin: %v\n", err)
	}

	authClient, err := client.App.Auth(context.Background())
	if err != nil {
		log.Fatalf("Error creating Auth client: %v\n", err)
	}

	decodedToken, err := authClient.VerifyIDToken(context.Background(), token)

	if err != nil {
		log.Fatalf("Error verifying ID token: %v\n", err)
	}

	fmt.Printf("Decoded Token: %+v\n", decodedToken.UID)
	fmt.Println("JWT Token:", token)

}

func signInWithEmailPassword(email, password string) (string, error) {
	apiKey := ""
	url := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=%s", apiKey)

	requestData := map[string]string{
		"email":             email,
		"password":          password,
		"returnSecureToken": "true",
	}

	jsonValue, _ := json.Marshal(requestData)
	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	var authResp firebaseAuthResponse
	if err := json.Unmarshal(body, &authResp); err != nil {
		return "", err
	}

	return authResp.IDToken, nil
}
