package utils

import (
	"log"
	"net/http"
	"os"
)

func GetEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func ExtractFormValue(w http.ResponseWriter, r *http.Request, field string) (string, bool) {
	err := r.ParseForm()
	if err != nil {
		log.Println("Error parsing form:", err)
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return "", true
	}

	extractedValue := r.Form.Get(field)
	if extractedValue == "" {
		http.Error(w, field+" is required", http.StatusBadRequest)
		return "", true
	}
	return extractedValue, false
}
