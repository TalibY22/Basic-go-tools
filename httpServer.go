package main

import (
	"fmt"
	"net/http"
	
)



func helloworld(w http.ResponseWriter,r *http.Request)  {
	fmt.Fprintf(w, "Hello %s\n", r.URL.Query().Get("name"))
}


func main()  {
	http.HandleFunc("/hello", helloworld)
	http.ListenAndServe(":8000", nil)
}