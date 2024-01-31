package main

import (
	"cl-inter/helpers"

	"github.com/gin-gonic/gin"
)

func View(c *gin.Context) {
	var data helpers.FormData
	c.Bind(&data)

	inputData, err := helpers.ConvertStrToMap(c, data.Data)
	if err == true {
		return
	}
	channel := make(chan map[string]interface{})

	go worker(channel)

	channel <- inputData

	response := helpers.ResponseStruct{
		Success: true,
		Message: "Message sent successfully.!",
	}

	c.JSON(200, response)
}

func main() {
	r := gin.Default()

	r.POST("", View)

	r.Use(helpers.CommonMiddleware())
	r.SetTrustedProxies([]string{"127.0.0.1"})
	r.Run(":8080")
}
