package initializers

import (
	firebase "firebase.google.com/go/v4"
	"golang.org/x/net/context"
	"google.golang.org/api/option"

	"log"
)

var Ctx = context.Background()

var FB *firebase.App

func ConnectToFirebase() {
	var err error
	conf := option.WithCredentialsFile("key/golang-notification-9d6a2-firebase-adminsdk-kcf5i-9480c6eaeb.json")
	FB, err = firebase.NewApp(Ctx, nil, conf)

	if err != nil {
		log.Printf("Error initializing app : %v", err)
	}
}
