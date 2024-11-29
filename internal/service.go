package internal

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"stulej-finder/internal/handlers"

	v1 "stulej-finder/internal/v1"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func Service() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load environment variables %v", err)
	}
	initDB()
	defer conn.Close()
	fmt.Println(os.Getenv("BOT_USERNAME"), os.Getenv("TWITCH_OAUTH"))
	client := twitch.NewClient(os.Getenv("BOT_USERNAME"), "oauth:"+os.Getenv("ACCESS_TOKEN"))
	var portString string = os.Getenv("PORT")

	go handlers.StartWsServer(os.Getenv("CLIENT_ID"), os.Getenv("ACCESS_TOKEN"), "")

	routes := v1.InitRoutes(*apiCfg)
	go func() {
		log.Printf("Server starting on port %v\n", portString)
		srv := &http.Server{
			Addr:    ":" + portString,
			Handler: routes,
		}
		log.Fatal(srv.ListenAndServe())
	}()

	monitorStreamer(client)
	client.OnPrivateMessage(func(msg twitch.PrivateMessage) {
		// Fetch the streamer_id based on the channel
		streamer, err := queries.GetStreamerByTwitchID(context.Background(), msg.RoomID)
		if err != nil {
			log.Printf("Error fetching streamer for channel %s: %v", msg.Channel, err)
			return
		}
		fmt.Println(streamer.ID, streamer.TwitchID)

		// Skip messages from StreamElements bot
		if msg.User.DisplayName == "StreamElements" {
			return
		}

		// Store the message for the specific streamer
		storeMessageForStreamer(msg, streamer.ID)
		log.Printf("[%s] %s: %s", msg.Channel, msg.User.DisplayName, msg.Message)
	})

	// Connect to Twitch
	err = client.Connect()
	if err != nil {
		log.Fatal(err)
	}
}

func monitorStreamer(client *twitch.Client) {
	go func() {
		for {
			streamers, err := queries.GetStreamers(context.Background())
			if err != nil {
				log.Printf("Error fetching streamers: %v", err)
				time.Sleep(1 * time.Minute)
				continue
			}

			for _, streamer := range streamers {
				client.Join(streamer.Username)
				log.Printf("Joined channel: %s", streamer.Username)
			}
			time.Sleep(5 * time.Minute)
		}
	}()
}
