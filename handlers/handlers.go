package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var Tmpl *template.Template

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		RenderError(w, http.StatusNotFound, "404 Not Found")
		return
	}
	if r.Method != http.MethodGet {
		RenderError(w, http.StatusMethodNotAllowed, "405 Method Not Allowed")
		return
	}

	// Safety check to ensure the template was initialized
	if Tmpl == nil {
		RenderError(w, http.StatusInternalServerError, "500 Internal Server Error: Template uninitialized")
		return
	}

	// Render index.html
	err := Tmpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		RenderError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/login/" {
		RenderError(w, http.StatusNotFound, "404 Not Found")
		return
	}
	if r.Method != http.MethodGet {
		RenderError(w, http.StatusMethodNotAllowed, "405 Method Not Allowed")
		return
	}

	// Safety check to ensure the template was initialized
	if Tmpl == nil {
		RenderError(w, http.StatusInternalServerError, "500 Internal Server Error: Template uninitialized")
		return
	}

	// Render index.html
	err := Tmpl.ExecuteTemplate(w, "login.html", nil)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		RenderError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/register/" {
		RenderError(w, http.StatusNotFound, "404 Not Found")
		return
	}
	if r.Method != http.MethodGet {
		RenderError(w, http.StatusMethodNotAllowed, "405 Method Not Allowed")
		return
	}

	// Safety check to ensure the template was initialized
	if Tmpl == nil {
		RenderError(w, http.StatusInternalServerError, "500 Internal Server Error: Template uninitialized")
		return
	}

	// Render index.html
	err := Tmpl.ExecuteTemplate(w, "register.html", nil)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		RenderError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

}

func ActivityHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/activity/" {
		RenderError(w, http.StatusNotFound, "404 Not Found")
		return
	}
	if r.Method != http.MethodGet {
		RenderError(w, http.StatusMethodNotAllowed, "405 Method Not Allowed")
		return
	}

	// Safety check to ensure the template was initialized
	if Tmpl == nil {
		RenderError(w, http.StatusInternalServerError, "500 Internal Server Error: Template uninitialized")
		return
	}

	// Render activity.html
	err := Tmpl.ExecuteTemplate(w, "activity.html", nil)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		RenderError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

}

func RenderError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, "%s", message)
}
