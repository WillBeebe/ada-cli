package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	authServerPort = 8080
	tokenFile      = "token23.json"
)

var (
	oauth2Config *oauth2.Config
)

func init() {
	oauth2Config = &oauth2.Config{
		ClientID:    "368747471698-eq1l90u2baqmu6gtm3fkg1k860kf5v64.apps.googleusercontent.com",
		RedirectURL: fmt.Sprintf("http://localhost:%d/callback", authServerPort),
		Scopes:      []string{"openid", "profile", "email"},
		Endpoint:    google.Endpoint,
	}
}

func Login() {
	// Check if we have a stored token
	token, err := loadToken()
	if err == nil {
		fmt.Println("Using existing token")
		// Use this token for API requests
		useToken(token)
		return
	}

	// Start the OAuth flow
	authURL := oauth2Config.AuthCodeURL("state", oauth2.AccessTypeOffline)

	// Start a web server to handle the callback
	http.HandleFunc("/callback", handleCallback)
	go http.ListenAndServe(fmt.Sprintf(":%d", authServerPort), nil)

	// Open the user's browser to the consent page
	err = openBrowser(authURL)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Waiting for OAuth callback...")
	// The handleCallback function will receive the token and save it
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	token, err := oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = saveToken(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Authentication successful! You can close this window.")
	log.Println("Token received and saved")

	// Use the token for API requests
	useToken(token)

	os.Exit(0)
}

func saveToken(token *oauth2.Token) error {
	data, err := json.Marshal(token)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(tokenFile, data, 0600)
}

func loadToken() (*oauth2.Token, error) {
	data, err := ioutil.ReadFile(tokenFile)
	if err != nil {
		return nil, err
	}
	var token oauth2.Token
	err = json.Unmarshal(data, &token)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func useToken(token *oauth2.Token) {
	// Use this token to make authenticated requests to your API
	client := oauth2Config.Client(context.Background(), token)
	resp, err := client.Get("http://localhost:8000/projects/")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("API response:", string(body))
}

func openBrowser(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	return err
}
