package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

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
	upgrader     = websocket.Upgrader{}
	clients      = make(map[*websocket.Conn]bool)
	clientsMutex sync.Mutex
)

func isTwitchStreamLive(clientID, accessToken, streamerName string) (bool, error) {
	// Twitch Helix API URL
	apiURL := "https://api.twitch.tv/helix/streams"

	// Prepare the request
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return false, err
	}

	// Add query parameters
	query := url.Values{}
	query.Add("user_login", streamerName)
	req.URL.RawQuery = query.Encode()

	// Set headers
	req.Header.Set("Client-ID", clientID)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// Parse the response
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("failed to fetch stream status: %s", resp.Status)
	}
	var streamResp StreamResponse
	if err := json.NewDecoder(resp.Body).Decode(&streamResp); err != nil {
		return false, err
	}

	// Check if the stream is live
	return len(streamResp.Data) > 0, nil
}

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

func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Websocket upgrade error: %v", err)
		return
	}
	clientsMutex.Lock()
	clients[conn] = true
	clientsMutex.Unlock()
	log.Println("New ws connection established")
	defer func() {
		clientsMutex.Lock()
		delete(clients, conn)
		clientsMutex.Unlock()
		conn.Close()
		log.Println("Ws conn closed")
	}()
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Ws read error: %v\n", err)
			break
		}
	}
}

func StartWsServer(clientID, accessToken, streamerName string) {
	go func() {
		for {
			isLive, err := isTwitchStreamLive(clientID, accessToken, streamerName)
			if err != nil {
				log.Printf("Error checking stream %v", err)
				time.Sleep(5 * time.Second)
				continue
			}
			message := fmt.Sprintf("The stream for %s is %s", streamerName, map[bool]string{true: "LIVE", false: "OFFLINE"}[isLive])
			log.Println(message)
			broadcastMessage(message)
			time.Sleep(5 * time.Second)
		}
	}()
}
