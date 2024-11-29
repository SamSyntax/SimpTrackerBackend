package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/gorilla/websocket"
)

type StreamResponse struct {
	Data []struct {
		ID          string `json:"id"`
		UserName    string `json:"user_name"`
		GameName    string `json:"game_name"`
		Title       string `json:"title"`
		ViewerCount int    `json:"viewer_count"`
		StartedAt   string `json:"started_at"`
	} `json:"data"`
}

var (
	allowedOrigins = map[string]bool{
		"https://simptracker.framed-designs.com": true,
	}
	upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		if strings.HasPrefix(origin, "http://localhost:") || strings.HasPrefix(origin, "chrome-extension://jlmpjdjjbgclbocgajdjefcidcncaied") {
			return true
		}
		return allowedOrigins[origin]
	}}
	clients      = make(map[*websocket.Conn]bool)
	clientsMutex sync.Mutex
	twitchClient *twitch.Client
)

// Initialize the Twitch IRC client
func initTwitchClient() {
	if twitchClient == nil {
		twitchClient = twitch.NewAnonymousClient()
	}
}

// Check if the Twitch stream is live
func isTwitchStreamLive(clientID, accessToken, streamerName string) (bool, error) {
	apiURL := "https://api.twitch.tv/helix/streams"
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %v", err)
	}
  fmt.Println("streamer:",streamerName)

	query := url.Values{}
	query.Add("user_login", streamerName)
	req.URL.RawQuery = query.Encode()

	req.Header.Set("Client-ID", clientID)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	log.Printf("Request URL: %s\n", req.URL.String())
	log.Printf("Authorization Header: %s\n", req.Header.Get("Authorization"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("failed to fetch stream status: %s", resp.Status)
	}

	var streamResp StreamResponse
	if err := json.NewDecoder(resp.Body).Decode(&streamResp); err != nil {
		return false, fmt.Errorf("failed to decode response: %v", err)
	}
	if strings.ToLower(os.Getenv("ENV")) == "dev" {
		return true, nil
	} else {
		return len(streamResp.Data) > 0, nil
	}
}

// Broadcast a message to all WebSocket clients
func broadcastMessage(message string) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Printf("Error writing to client: %v\n", err)
			client.Close()
			delete(clients, client)
		}
	}
}

// WebSocket handler for client connections
func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	defer func() {
		conn.Close()
		log.Println("WebSocket connection closed")
	}()

	clientsMutex.Lock()
	clients[conn] = true
	clientsMutex.Unlock()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v\n", err)
			break
		}
	}

	// Clean up the client connection when the loop exits
	clientsMutex.Lock()
	delete(clients, conn)
	clientsMutex.Unlock()
}

// Monitor the streamer's status and handle chat connection
func monitorStreamer(clientID, accessToken, streamerName string) {
	initTwitchClient()
	streamerOnline := false

	for {
		isLive, err := isTwitchStreamLive(clientID, accessToken, streamerName)
		if err != nil {
			log.Printf("Error checking stream status: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		if isLive && !streamerOnline {
			// Connect to chat when streamer goes live
			log.Println("Streamer is live. Connecting to chat...")
			twitchClient.Join(streamerName)
			twitchClient.OnPrivateMessage(func(msg twitch.PrivateMessage) {
				log.Printf("[%s] %s: %s", msg.Channel, msg.User.DisplayName, msg.Message)
				broadcastMessage(fmt.Sprintf("[%s] %s: %s", msg.Channel, msg.User.DisplayName, msg.Message))
			})

			go func() {
				if err := twitchClient.Connect(); err != nil {
					log.Printf("Error connecting to Twitch chat: %v", err)
				}
			}()
			streamerOnline = true
			broadcastMessage("LIVE")
		} else if !isLive && streamerOnline {
			// Disconnect from chat when streamer goes offline
			log.Println("Streamer is offline. Disconnecting from chat...")
			twitchClient.Depart(streamerName)
			err = twitchClient.Disconnect()
			if err != nil {
				log.Printf("Failed to disconnect twitch client: %v", err)
			}
			streamerOnline = false
			broadcastMessage("OFFLINE")
		} else {
			// Send periodic status update
			status := map[bool]string{true: "LIVE", false: "OFFLINE"}[isLive]
			broadcastMessage(status)
		}

		time.Sleep(10 * time.Second) // Poll every 10 seconds
	}
}

// Start WebSocket server and monitor streamer
func StartWsServer(clientID, accessToken, streamerName string) {
	go monitorStreamer(clientID, accessToken, streamerName)
	log.Println("WebSocket server started")
}
