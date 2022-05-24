package poststore

type Config struct {
	Id      string            `json: "id"`
	Entries map[string]string `json:"entries"`
	Version string            `json:"version"`
}

type Group struct {
	Id      string
	Configs []Config
	Version string `json: "version"`
}
