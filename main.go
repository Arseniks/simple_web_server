package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var arseniy = &person{
	Name: "Arseniy",
	Age:  18,
}

func personHandler(writer http.ResponseWriter, request *http.Request) {

	switch request.Method {
	case "GET":
		response, err := json.Marshal(arseniy)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
		_, err = writer.Write(response)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
	case "POST":
		d := json.NewDecoder(request.Body)
		newPerson := &person{}
		err := d.Decode(newPerson)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
		arseniy = newPerson
	default:
		writer.WriteHeader(http.StatusMethodNotAllowed)
		_, err := fmt.Fprintf(writer, "I can't do that.")
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
	}
}

func main() {
	http.HandleFunc("/person/", personHandler)

	log.Println("Staring server on: http://localhost:8080 ")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}
