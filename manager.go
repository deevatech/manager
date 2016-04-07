package main

import (
	"fmt"
	"github.com/deevatech/manager/runner"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

func init() {
	log.Println("Deeva Manager!")
}

func main() {
	router := gin.Default()
	router.POST("/run", handleRunRequest)

	port := os.Getenv("DEEVA_MANAGER_PORT")
	if len(port) == 0 {
		port = "9090"
	}

	log.Printf("Starting in %s mode on port %s\n", gin.Mode(), port)
	host := fmt.Sprintf(":%s", port)
	router.Run(host)
}

func handleRunRequest(c *gin.Context) {
	if err := runner.Run(); err != nil {
		log.Println(err)
	}

	c.JSON(http.StatusOK, gin.H{})
}
