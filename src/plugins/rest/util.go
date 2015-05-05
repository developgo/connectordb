package rest

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

//ErrInvalidName is thrown when the name is bad
var ErrInvalidName = errors.New("The given name did not pass sanitation.")

//OK is a simplifying function that returns success
func OK(writer http.ResponseWriter) error {
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("ok"))
	log.Println("OK")
	return nil
}

//JSONWriter writes the given data as http
func JSONWriter(writer http.ResponseWriter, data interface{}, err error) error {
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return err
	}

	res, err := json.Marshal(data)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return err
	}
	writer.WriteHeader(http.StatusOK)
	writer.Write(res)
	log.Printf("Success: %s\n", string(res))
	return nil
}

//UnmarshalRequest unmarshals the input data to the given interface
func UnmarshalRequest(request *http.Request, unmarshalTo interface{}) error {
	defer request.Body.Close()

	//Limit requests to 10MB
	data, err := ioutil.ReadAll(io.LimitReader(request.Body, 10000000))
	if err != nil {
		return err
	}

	return json.Unmarshal(data, unmarshalTo)
}

//ValidName sanitizes names so that only valid ones are added
func ValidName(n string, err error) error {
	if err != nil {
		return err
	}
	if strings.Contains(n, "/") || strings.Contains(n, "\\") || strings.Contains(n, " ") || strings.Contains(n, "?") {
		return ErrInvalidName
	}
	if n == "ls" || n == "this" {
		return ErrInvalidName
	}
	return nil
}
