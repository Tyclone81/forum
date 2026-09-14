package main

import (
	"fmt"
	"forum/handlers"
	"html/template"
	"log"
	"net/http"
)

func main() {
	var err error

	//FRONT-END LOAD
	handlers.Tmpl, err = template.ParseFiles(
		"templates/partials.html",
		"templates/index.html",
		"templates/login.html",
		"templates/register.html",
		"templates/activity.html",
	)
	if err != nil {
		log.Fatalf("Initialization Error: Failed to parse templates: %v", err)
	}

	cssFs := http.FileServer(http.Dir("./css/"))
	http.Handle("/css/", http.StripPrefix("/css/", cssFs))

	jsFs := http.FileServer(http.Dir("./js/"))
	http.Handle("/js/", http.StripPrefix("/js/", jsFs))

	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/register/", handlers.RegisterHandler)
	http.HandleFunc("/login/", handlers.LoginHandler)
	http.HandleFunc("/activity/", handlers.ActivityHandler)

	fmt.Println("Server executing stably at http://localhost:8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Fatal: Server failed to bind to port: %v", err)

	}
}
