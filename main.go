package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Article struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `json:"title"`
	Slug      string         `gorm:"unique_index" json:"slug"`
	Desc      string         `gorm:"type:text;" json:"desc"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

var DB *gorm.DB

func main() {
	var err error

	err = godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbTimeZone := os.Getenv("DB_TIMEZONE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s", dbHost, dbUser, dbPass, dbName, dbPort, dbTimeZone)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	sqlDB, _ := DB.DB()
	defer sqlDB.Close()

	// Migrate the schema
	DB.AutoMigrate(&Article{})

	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		articles := v1.Group("/articles")
		{
			articles.GET("/", getArticles)
			articles.GET("/:slug", getArticle)
			articles.POST("/", postArticle)
		}
	}

	router.Run()
}

func getArticles(c *gin.Context) {
	items := []Article{}
	DB.Find(&items)

	c.JSON(200, gin.H{
		"status": "berhasil ke halaman home",
		"data":   items,
	})
}

func getArticle(c *gin.Context) {
	slug := c.Param("slug")

	var item Article

	result := DB.First(&item, "slug = ?", slug)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(404, gin.H{"status": "error", "message": "record not found"})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{
		"status": "berhasil",
		"data":   item,
	})
}

func postArticle(c *gin.Context) {
	item := Article{
		Title: c.PostForm("title"),
		Desc:  c.PostForm("desc"),
		Slug:  slug.Make(c.PostForm("title")),
	}

	DB.Create(&item)

	c.JSON(200, gin.H{
		"status": "berhasil ngepost",
		"data":   item,
	})
}
