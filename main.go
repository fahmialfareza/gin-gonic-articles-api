package main

import "github.com/gin-gonic/gin"

func main() {
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		articles := v1.Group("/articles")
		{
			articles.GET("/", getHome)
			articles.GET("/:title", getArticle)
			articles.POST("/", postArticle)
		}
	}

	router.Run()
}

func getHome(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "berhasil",
		"message": "Berhasil akses home",
	})
}

func getArticle(c *gin.Context) {
	title := c.Param("title")

	c.JSON(200, gin.H{
		"status":  "berhasil",
		"message": title,
	})
}

func postArticle(c *gin.Context) {
	title := c.PostForm("title")
	desc := c.PostForm("desc")

	c.JSON(200, gin.H{
		"status":  "berhasil",
		"message": title,
		"desc":    desc,
	})
}
