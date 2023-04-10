package emailhider

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"gopkg.in/ezzarghili/recaptcha-go.v4"
)

func init() {
	functions.HTTP("handleRequest", handleRequest)
}

type Request struct {
	Token string
	Site  string
}

type Response struct {
	Email string `json:"email"`
}

func handleRequest(w http.ResponseWriter, r *http.Request) {

	// Set CORS headers for the preflight request
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	fmt.Println(r.Header.Get("origin"))

	if r.Method != "POST" {
		log.Println("Error: " + r.Method + " not allowed.")
		http.Error(w, "Only POST requests are accepted.", http.StatusMethodNotAllowed)
		return

	}

	captcha, _ := recaptcha.NewReCAPTCHA(os.Getenv("RECAPTCHA_SECRET"), recaptcha.V2, 10*time.Second) // for v2 API get your secret from https://www.google.com/recaptcha/admin

	var body Request

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		// body should be json
		log.Println("Error: " + err.Error())

		http.Error(w, "Bad request body.", http.StatusBadRequest)
		return
	}

	if body.Token != "" {
		if err := captcha.Verify(body.Token); err != nil {
			// recaptcha validation failed
			log.Println("Error: " + err.Error())

			http.Error(w, "Validation failed.", http.StatusForbidden)
			return
		}
	} else {
		// recaptcha token not provided
		log.Println("Error: token not provided")

		http.Error(w, "No recaptcha token provided.", http.StatusBadRequest)
		return
	}

	if body.Site != "" {
		email := os.Getenv(body.Site)
		if email == "" {
			// site provided was not found
			log.Println("Error: " + body.Site + " site not available")

			http.Error(w, "Validation failed.", http.StatusForbidden)
			return
		}

		response := Response{Email: email}
		fmt.Println(body.Site + ": " + response.Email)

		// all is good, return the email address
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	} else {
		// site not provided
		log.Println("Error: site not provided")

		http.Error(w, "No site provided.", http.StatusBadRequest)
		return
	}
}
