package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func TrimJson(jsonBytes []byte) []byte {
	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, jsonBytes); err != nil {
		fmt.Println(err)
	}

	return buffer.Bytes()
}

func GetDomainWithHost(host string) string {
	domain := strings.Split(host, ":")
	domain = strings.Split(domain[0], ".")

	return domain[0]
}

func Handler404(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)

	err := json.NewEncoder(w).Encode(map[string]interface{}{
		"status": http.StatusText(http.StatusNotFound),
	})

	if err != nil {
		log.Printf("[ERROR] Error response error: %v", err)
		return
	}

	log.Printf("[DEBUG] [%s] 404 %s", r.Method, r.URL.String())
}
