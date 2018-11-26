package config

type Binary struct {
	Root  string            `json:"root"`
	Files map[string]string `json:"files"`
}
