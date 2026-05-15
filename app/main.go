package main

import (
	"encoding/json"
	"net/http"
)

// Response define a estrutura do JSON de retorno
type Response struct {
	Message string `json:"message"`
}

// helloHandler manipula as requisições HTTP na rota principal
func helloHandler(w http.ResponseWriter, r *http.Request) {
	// Define o cabeçalho para retornar o formato JSON
	w.Header().Set("Content-Type", "application/json")

	// Define o status HTTP 200 OK
	w.WriteHeader(http.StatusOK)

	// Instancia e envia a resposta formatada
	response := Response{Message: "Hello, World!"}
	json.NewEncoder(w).Encode(response)
}

func readyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func main() {
	// Vincula a rota raiz "/" à função do manipulador
	http.HandleFunc("/ready", readyHandler)

	// Inicia o servidor HTTP local na porta 8080
	http.ListenAndServe(":6969", nil)
}
