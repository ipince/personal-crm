package main

import (
	"context"
	"easyrm/people"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2/google"
	peoplev1 "google.golang.org/api/people/v1"
)

func main() {
	ctx := context.Background()
	b, err := os.ReadFile("oauth_desktop_client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, peoplev1.ContactsReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	srv := people.Client(ctx, config)

	people.List(srv, 10)
	test, err := people.Get(srv, people.TestPersonID)
	if err != nil {
		fmt.Println(err)
	}

	bytes, err := test.MarshalJSON()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(bytes))

	err = people.SetBirthdate(srv, people.TestPersonID, 1, 2, nil)
	if err != nil {
		fmt.Println(err)
	}
}
