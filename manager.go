package main

import (
	"fmt"
	"github.com/deevatech/manager/models/tests"
	"github.com/deevatech/manager/runner"
	. "github.com/deevatech/manager/types"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strconv"
)

func init() {
	log.Println("Deeva Manager!")
}

func main() {
	router := gin.Default()
	router.POST("/run", handleRunRequest)
	router.GET("/tests/:id", handleTestLookupRequest)
	router.POST("/tests/:id/submit", handleTestSubmitRequest)

	port := os.Getenv("DEEVA_MANAGER_PORT")
	if len(port) == 0 {
		port = "8080"
	}

	log.Printf("Starting in %s mode on port %s\n", gin.Mode(), port)
	host := fmt.Sprintf(":%s", port)
	router.Run(host)
}

func handleRunRequest(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")

	var run RunParams
	if errParams := c.BindJSON(&run); errParams == nil {
		if result, errRun := runner.Run(run); errRun == nil {
			c.JSON(http.StatusOK, result)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": errRun,
			})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errParams,
		})
	}
}

func handleTestLookupRequest(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	test := tests.FindById(id)
	c.JSON(http.StatusOK, test)
}

func handleTestSubmitRequest(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	// lookup test
	test := tests.FindById(id)

	var submit TestSubmitParams
	if errParams := c.BindJSON(&submit); errParams == nil {
		log.Printf("TestSubmitParams: %#v", submit)
		run := RunParams{
			Language: test.Language,
			Source:   submit.Code,
			Spec:     test.Spec,
		}
		log.Printf("RunParams: %#v", run)
		if result, errRun := runner.Run(run); errRun == nil {
			c.JSON(http.StatusOK, result)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": errRun,
			})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errParams,
		})
	}
}
