package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

type PredictionResult struct {
	Prediction string `json:"prediction"`
}

func main() {
	router := NewRouter()
	log.Fatal(http.ListenAndServe(":8000", router))
}

func NewRouter() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("/predict", PredictHandler)
	return router
}

type UserInput struct {
	Pregnancies              int     `json:"Pregnancies"`
	Glucose                  int     `json:"Glucose"`
	BloodPressure            int     `json:"BloodPressure"`
	SkinThickness            int     `json:"SkinThickness"`
	Insulin                  int     `json:"Insulin"`
	BMI                      float64 `json:"BMI"`
	DiabetesPedigreeFunction float64 `json:"DiabetesPedigreeFunction"`
	Age                      int     `json:"Age"`
}

func PredictHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	fmt.Println("received request")

	decoder := json.NewDecoder(r.Body)
	var userInput UserInput
	err := decoder.Decode(&userInput)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid JSON format")
		return
	}

	pythonScript := "dia.py"
	jsonInput, err := json.Marshal(userInput)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to encode JSON")
		return
	}

	cmd := exec.Command("python", pythonScript, string(jsonInput))
	output, err := cmd.Output()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to execute Python script")
		return
	}

	predictionResult := PredictionResult{
		Prediction: string(output),
	}

	responseJSON, err := json.Marshal(predictionResult)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to encode JSON")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}
