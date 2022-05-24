package database

import (
	"fmt"

	"github.com/google/uuid"
)

const (
	posts = "posts/%s"
	all   = "posts"
)

func generateKey() (string, string) {
	id := uuid.New().String()
	return fmt.Sprintf(posts, id), id
}

func constructKey(id string) string {
	return fmt.Sprintf(posts, id)
}
