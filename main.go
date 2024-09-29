package main

import (
	"fmt"
	"library-system/conf"
	"library-system/database"
	"library-system/handler"
	"log"

	"github.com/gin-gonic/gin"
)

const configFilePath = "config.json"

func main() {
	config := conf.LoadConfiguration(configFilePath)
	serverURL := fmt.Sprintf("localhost:%d", config.ServicePort)
	dbConn, dberr := database.NewDBConnection(config)
	if dberr != nil {
		log.Fatalf("Failed to connect to Database %v", dberr)
	}
	handler := handler.NewCommandHandler(dbConn)
	r := gin.Default()
	v1 := r.Group("")
	{
		v1.GET("/books", handler.ListBooks)

		v1.POST("/books", handler.AddBook)
		v1.POST("/books/:id/borrow", handler.BorrowBook)
		v1.POST("/books/:id/return", handler.ReturnBook)
	}
	r.Run(serverURL)

}
