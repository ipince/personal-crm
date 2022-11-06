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

	validateAndNormalizeAll(srv)

	testChanges(srv)
}

func fetchAll(srv *peoplev1.Service) []*peoplev1.Person {
	all, err := people.ListAll(srv)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Fetched %d connections\n", len(all))
	return all
}

func testChanges(srv *peoplev1.Service) {

	//test, err := people.Get(srv, people.TestPersonID)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//people.Print(test)
}

func validateAndNormalizeAll(srv *peoplev1.Service) {
	all := fetchAll(srv)
	for _, p := range all {
		err := people.Normalize(srv, p)
		if err != nil {
			fmt.Println(err)
		}
		people.Validate(p)
	}
}
