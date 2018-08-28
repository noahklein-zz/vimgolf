package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	_ "github.com/go-sql-driver/mysql"

	"github.com/noahklein/vimgolf/storage"
)

var (
	port   = flag.Int("port", 3000, "The port to serve on")
	dbAddr = flag.String("dbAddr", "", "The DB address")
)

func main() {
	flag.Parse()

	storage.Init("mysql", *dbAddr)

	r := gin.Default()

	{
		r.LoadHTMLGlob("./client/**/*.html")
		r.GET("/", indexPage)
		r.GET("/challenge/:id", challengePage)
	}

	api := r.Group("/api")
	{
		api.POST("challenge/:id/answer", answer)
	}

	fmt.Printf("Serving... http://localhost:%d", *port)
	r.Run(fmt.Sprintf(":%d", *port))
}

func indexPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func challengeID(c *gin.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

func challengePage(c *gin.Context) {
	id, err := challengeID(c)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
	}

	challenge, err := storage.ChallengeByID(id)
	c.HTML(http.StatusOK, "challenge.html", gin.H{
		"start":  challenge.Start,
		"target": challenge.Target,
	})
}

func answer(c *gin.Context) {
	type answerRequest struct {
		Cmd  string `json:"cmd" binding:"required"`
		User string `json:"user" binding:"required"`
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	id, err := challengeID(c)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}
	var req answerRequest
	c.BindQuery(&req)

	challenge, err := storage.ChallengeByID(id)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}
	solution, err := challenge.Answer(ctx, req.Cmd, req.User)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, solution)
}
