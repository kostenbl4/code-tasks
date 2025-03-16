package domain

type Task struct {
	Translator string `json:"translator"`
	Code       string `json:"code"`
	UUID       string `json:"uuid"`
	Status     string `json:"status"`
	Result     string `json:"result"`
	Stdout     string `json:"stdout"`
	Stderr     string `json:"stderr"`
}