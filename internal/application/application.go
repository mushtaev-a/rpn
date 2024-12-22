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

var BadRequestErrors []error = []error{
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
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid JSON request body", http.StatusBadRequest)
		return
	}

	if reqBody.Expression == "" {
		http.Error(w, "Empty expression", http.StatusBadRequest)
		return
	}

	result, err := rpn.Calc(strings.TrimSpace(reqBody.Expression))

	if err != nil {
		statusCode := http.StatusInternalServerError

		for _, e := range BadRequestErrors {
			if errors.Is(err, e) {
				statusCode = http.StatusBadRequest
				http.Error(w, err.Error(), statusCode)
				return
			}
		}

		http.Error(w, err.Error(), statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float64{"result": result})
}

func (a *Application) Run() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", HandleCalculation)

	log.Printf("Server is running on :%s port", a.config.Port)

	log.Fatal(http.ListenAndServe(":"+a.config.Port, mux))
}
