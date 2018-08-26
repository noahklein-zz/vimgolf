package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/noahklein/vimgolf/vim"
)

var (
	port = flag.Int("port", 3000, "The port to serve on")
)

func main() {
	flag.Parse()

	r := gin.Default()
	api := r.Group("/api")
	{
		api.POST("problem/:problem/answer", answer)

		fmt.Printf("Serving on port %d...", *port)
		r.Run(fmt.Sprintf(":%d", *port))
	}
}

type AnswerRequest struct {
	Cmd string `json:"cmd" binding:"required"`
}

func answer(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var req AnswerRequest
	c.BindQuery(&req)
	problem := c.Param("problem")

	got, err := vim.AttemptChallenge(ctx, problem, req.Cmd)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	if got == "" {

	}
}
