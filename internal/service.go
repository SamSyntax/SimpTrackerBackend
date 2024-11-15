package internal

import (
	"log"
	"net/http"
	"os"

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
	client := twitch.NewAnonymousClient()
	var portString string = os.Getenv("PORT")

	go func() {
		routes := v1.InitRoutes(*apiCfg)
		log.Printf("Server starting on port %v\n", portString)
		srv := &http.Server{
			Addr:    ":" + portString,
			Handler: routes,
		}
		log.Fatal(srv.ListenAndServe())
	}()
	client.Join("karolynkaa")

	client.OnPrivateMessage(func(msg twitch.PrivateMessage) {
		if msg.User.DisplayName == "StreamElements" {
			return
		}
		storeMessage(msg)
		log.Printf("%s: %s", msg.User.DisplayName, msg.Message)
	})

	err = client.Connect()
	if err != nil {
		log.Fatal(err)
	}
}
