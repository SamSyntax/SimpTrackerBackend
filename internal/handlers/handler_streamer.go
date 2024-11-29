package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"stulej-finder/internal/db"
)

type TwitchTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

func (apiCfg *ApiConfig) HandlerAuthRedirect(w http.ResponseWriter, r *http.Request) {
	clientID := os.Getenv("CLIENT_ID")
	redirectURI := "http://localhost:8080/v1/auth/callback" // Ensure this matches your Twitch app's redirect URI
	scopes := "chat:read+chat:edit"

	authURL := fmt.Sprintf(
		"https://id.twitch.tv/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=%s",
		clientID, redirectURI, scopes,
	)

	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func (apiCfg *ApiConfig) HandlerAuthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		log.Println("Missing code in callback request")
		http.Error(w, "Missing code in callback request", http.StatusBadRequest)
		return
	}

	// Exchange the code for tokens
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	redirectURI := "http://localhost:8080/v1/auth/callback"

	tokenURL := "https://id.twitch.tv/oauth2/token"
	reqBody := fmt.Sprintf(
		"client_id=%s&client_secret=%s&code=%s&grant_type=authorization_code&redirect_uri=%s",
		clientID, clientSecret, code, redirectURI,
	)

	resp, err := http.Post(tokenURL, "application/x-www-form-urlencoded", strings.NewReader(reqBody))
	if err != nil {
		log.Printf("Failed to exchange code for tokens: %v\n", err)
		http.Error(w, "Failed to exchange code for tokens", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Token exchange failed: %s\n", string(body))
		http.Error(w, "Failed to exchange code for tokens", http.StatusInternalServerError)
		return
	}

	body, _ := io.ReadAll(resp.Body)
	var tokenResp TwitchTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		log.Printf("Failed to parse token response: %v\n", err)
		http.Error(w, "Failed to parse token response", http.StatusInternalServerError)
		return
	}

	// Fetch streamer details
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.twitch.tv/helix/users", nil)
	req.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)
	req.Header.Set("Client-ID", clientID)

	streamerResp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to fetch streamer details: %v\n", err)
		http.Error(w, "Failed to fetch streamer details", http.StatusInternalServerError)
		return
	}
	defer streamerResp.Body.Close()

	if streamerResp.StatusCode != http.StatusOK {
		streamerBody, _ := io.ReadAll(streamerResp.Body)
		log.Printf("Streamer details request failed: %s\n", string(streamerBody))
		http.Error(w, "Failed to fetch streamer details", http.StatusInternalServerError)
		return
	}

	streamerBody, _ := io.ReadAll(streamerResp.Body)
	var streamerData struct {
		Data []struct {
			ID       string `json:"id"`
			Login    string `json:"login"`
			Username string `json:"display_name"`
		} `json:"data"`
	}
	if err := json.Unmarshal(streamerBody, &streamerData); err != nil {
		log.Printf("Failed to parse streamer data: %v\n", err)
		http.Error(w, "Failed to parse streamer data", http.StatusInternalServerError)
		return
	}

	if len(streamerData.Data) == 0 {
		log.Println("No streamer data returned")
		http.Error(w, "No streamer data found", http.StatusBadRequest)
		return
	}

	streamer := streamerData.Data[0]

	// Save streamer to database
	expiresAt := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	_, err = apiCfg.DB.UpsertStreamer(r.Context(), db.UpsertStreamerParams{
		TwitchID:     streamer.ID,
		Username:     streamer.Username,
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		ExpiresAt:    expiresAt,
	})
	if err != nil {
		log.Printf("Failed to save streamer to database: %v\n", err)
		http.Error(w, "Failed to save streamer to database", http.StatusInternalServerError)
		return
	}

	log.Printf("Streamer %s authenticated successfully!\n", streamer.Username)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Streamer %s authenticated successfully!", streamer.Username)
}
