package main

import (
	"net/http"

	banzuke "github.com/Doarakko/otoko-banzuke/internal/banzuke"
	commend "github.com/Doarakko/otoko-banzuke/internal/commend"
	search "github.com/Doarakko/otoko-banzuke/internal/search"
	myyoutube "github.com/Doarakko/otoko-banzuke/pkg/youtube"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {

	// err := godotenv.Load("../.env")
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	router := gin.Default()
	router.Static("/web/static", "./web/static")
	router.LoadHTMLGlob("templates/*")

	// 番付
	rankComments := banzuke.SelectRankComments()
	todayComments := banzuke.SelectTodayComments()
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"rankComments":  rankComments,
			"todayComments": todayComments,
		})
	})

	// 漢を推薦する
	channels := []myyoutube.Channel{}
	router.GET("/commend", func(c *gin.Context) {
		c.HTML(http.StatusOK, "commend/index.tmpl", gin.H{
			"channels": channels,
		})
	})

	router.POST("/commend", func(c *gin.Context) {
		query := c.PostForm("query")
		channelID := c.PostForm("channel_id")

		if query != "" {
			channels = commend.SearchChannels(query)
			for i := range channels {
				channels[i].SetDetailInfo()
			}
		} else if channelID != "" {
			commend.InsertChannel(channelID)
		}
		c.Redirect(302, "/commend")
	})

	// 漢を探す
	comments := []myyoutube.Comment{}
	router.GET("/search", func(c *gin.Context) {
		c.HTML(http.StatusOK, "search/index.tmpl", gin.H{
			"comments": comments,
		})
	})

	router.POST("/search", func(c *gin.Context) {
		query := c.PostForm("query")
		comments = search.SearchOtoko(query)

		c.Redirect(302, "/search")
	})

	// local
	// router.Run(":8080")
	router.Run()
}