<<<<<<< Updated upstream
package database

import (
	"fmt"

	"github.com/google/uuid"
)

const (
	group          = "groups/%s/%s/%s/"
	groupVersion   = "groups/%s/%s"
	config         = "configs/%s/%s"
	allGroups      = "groups"
	allConfigs     = "configs"
	idempotencyKey = "request/%s"
)

func generateKey(newId string, version string, label string) (string, string) {
	id := uuid.New().String()
	if label == "" {
		if newId != "" {
			return fmt.Sprintf(config, newId, version), newId
		}
		return fmt.Sprintf(config, id, version), id

	} else {
		return fmt.Sprintf(group, newId, version, label), newId
	}
}

func constructKey(id string, version string, label string) string {
	if version == "" && label == "" {
		return fmt.Sprintf(idempotencyKey, id)
	} else if label == "" {
		return fmt.Sprintf(config, id, version)
	} else if label != "" {
		return fmt.Sprintf(group, id, version, label)
	}
	return ""
}

func constructConfigKey(id string, version string) string {
	return fmt.Sprintf(config, id, version)
}

func constructGroupKey(id string, version string) string {
	return fmt.Sprintf(groupVersion, id, version)
}
=======
package database

import (
	"fmt"

	"github.com/google/uuid"
)

const (
	group          = "groups/%s/%s/%s/"
	groupVersion   = "groups/%s/%s"
	config         = "configs/%s/%s"
	allGroups      = "groups"
	allConfigs     = "configs"
	idempotencyKey = "request/%s"
)

func generateKey(newId string, version string, label string) (string, string) {
	id := uuid.New().String()
	if label == "" {
		if newId != "" {
			return fmt.Sprintf(config, newId, version), newId
		}
		return fmt.Sprintf(config, id, version), id

	} else {
		return fmt.Sprintf(group, newId, version, label), newId
	}
}

func constructKey(id string, version string, label string) string {
	if version == "" && label == "" {
		return fmt.Sprintf(idempotencyKey, id)
	} else if label == "" {
		return fmt.Sprintf(config, id, version)
	} else {
		return fmt.Sprintf(group, id, version, label)
	}

}

func constructConfigKey(id string, version string) string {
	return fmt.Sprintf(config, id, version)
}

func constructGroupKey(id string, version string) string {
	return fmt.Sprintf(groupVersion, id, version)
}
>>>>>>> Stashed changes
