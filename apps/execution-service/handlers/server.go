package handlers

import (
	"cloud.google.com/go/firestore"
)

type Service struct {
	Client *firestore.Client
}
