package database

import (
	"fmt"

	"github.com/google/uuid"
)

const (
	group      = "groups/%s/%s/%s"
	config     = "configs/%s/%s"
	allGroups  = "groups"
	allConfigs = "configs"
)

func generateKey(version string, label string) (string, string) {
	id := uuid.New().String()
	if label == "" {
		return fmt.Sprintf(config, id, version), id

	} else {
		return fmt.Sprintf(group, id, version, label), id
	}
}

func constructKey(id string, version string, label string) string {
	if label == "" {
		return fmt.Sprintf(config, id, version)
	} else {
		return fmt.Sprintf(group, id, version, label)
	}

}
