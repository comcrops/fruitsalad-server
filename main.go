package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type MyStruct struct {
	Beidl int8
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}


func helloHandler(w http.ResponseWriter, req *http.Request) {
	myStruct := &MyStruct {
		Beidl: 3,
	}
	jsonData, err := json.Marshal(myStruct)

	if err != nil {
		log.Fatalf("SAMC")
	}

	io.WriteString(w, string(jsonData[:]))
}

