package types


type ProcessTask struct {
	Translator string `json:"translator"`
	Code       string `json:"code"`
	UUID       string `json:"uuid"`
	Status     string `json:"status"`
	Stdout     string `json:"stdout"`
	Stderr     string `json:"stderr"`
}
