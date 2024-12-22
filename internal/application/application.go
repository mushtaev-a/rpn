package application

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/mushtaev-a/rpn/pkg/rpn"
)

type Config struct {
	Port string
}

func ConfigFromEnv() *Config {
	config := new(Config)
	config.Port = os.Getenv("PORT")

	if config.Port == "" {
		config.Port = "8080"
	}

	return config
}

type Application struct {
	config *Config
}

func New() *Application {
	return &Application{
		config: ConfigFromEnv(),
	}
}

var UnprocessableRequestErrors []error = []error{
	rpn.ErrDividingByZero,
	rpn.ErrDuplicateOpertaionsSigns,
	rpn.ErrOpertaionsSigns,
	rpn.ErrExpressionStringEmpty,
	rpn.ErrExpressionStringParetheses,
}

func HandleCalculation(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Expression string `json:"expression"`
	}

	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON request body"})
		return
	}

	if reqBody.Expression == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Empty expression"})
		return
	}

	result, err := rpn.Calc(strings.TrimSpace(reqBody.Expression))

	if err != nil {
		statusCode := http.StatusInternalServerError
		w.Header().Set("Content-Type", "application/json")

		for _, e := range UnprocessableRequestErrors {
			if errors.Is(err, e) {
				statusCode = http.StatusUnprocessableEntity
				w.WriteHeader(statusCode)
				json.NewEncoder(w).Encode(map[string]string{"error": "Expression is not valid"})
				return
			}
		}

		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float64{"result": result})
}

func (a *Application) Run() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/calculate", HandleCalculation)

	log.Printf("Server is running on :%s port", a.config.Port)

	log.Fatal(http.ListenAndServe(":"+a.config.Port, mux))
}
