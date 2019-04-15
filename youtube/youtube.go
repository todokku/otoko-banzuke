package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

func newGormConnect() *gorm.DB {
	db, err := gorm.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		panic(err.Error())
	}
	return db
}

func newYoutubeService() *youtube.Service {
	client := &http.Client{
		Transport: &transport.APIKey{Key: os.Getenv("YOTUBE_API_KEY")},
	}

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Unable to create YouTube service: %v", err)
	}

	return service
}

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	channel := channel{ChannelID: "UCxbY38ReXW3LbaviWUE4omg"}

	for _, video := range channel.selectVideos() {
		for _, comment := range video.getComments() {
			if comment.checkOtoko() {
				println(comment.TextDisplay)
			}
		}
	}
}