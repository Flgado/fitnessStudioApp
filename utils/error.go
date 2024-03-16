package utils

import "encoding/json"

type Error struct {
	Code        int
	err         error
	messages    map[string]string
	details     string
	suggestions string
}

func (e Error) Error() string {
	jsonData, _ := json.Marshal(e)
	return string(jsonData)
}

func E(code int, e error, message map[string]string, details string, suggestions string) Error {
	return Error{
		Code:        code,
		err:         e,
		messages:    message,
		details:     details,
		suggestions: suggestions,
	}
}

func (e Error) Message() map[string]string {
	return e.messages
}

func (e Error) StatusCode() int {
	return e.Code
}

func (e Error) GetErrorDetais() string {
	return e.details
}
func (e Error) GetSuguestions() string {
	return e.suggestions
}
