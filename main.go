package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main(){
	http.HandleFunc("/", func(rw http.ResponseWriter, r*http.Request) {
		log.Printf("Hello World")
		d, err := ioutil.ReadAll(r.Body)

		if err != nil{
			http.Error(rw, "Ooops", http.StatusBadRequest)
			return
		}
		fmt.Fprintf(rw, "Hello %s\n", d)
	})

	http.ListenAndServe(":9090", nil)
}