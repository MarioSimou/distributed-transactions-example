package main

import (
	"fmt"
	"log"
	"net/http"
)

func ping(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "pong\n")
}

func main(){
	var port = fmt.Sprintf(":%d", 3000)
	
	http.HandleFunc("/ping", ping)

	fmt.Printf("The app listens on port %s\n", port)
	log.Fatalln(http.ListenAndServe(port, nil))
}