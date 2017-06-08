/*
This application implements a bare-bones goCMS plugin. It includes:
- basic Gin routes
- A Gin route group
- Go templates
- Gin middleware

Check the authorization API using the Intellij REST tool to manually:
1. POST {"email":"admin@gocms.io", "password":"password"} to /api/login
2. Copy the resulting x-auth-token into the GET request header
3. GET any supported path (e.g., /api/private)
See /docs/goCMS for all the routes supported by goCMS.

*/
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	// Allow the calling program to set the port on which this service runs; Tashi has generalized this feature.
	port := flag.Int("port", 8080, "allow goCMS to set the port via this commandline option...")
	flag.Parse()

	// Build the basic Gin router
	router := gin.Default()
	// Load a simple template
	router.LoadHTMLGlob(filepath.Join(os.Getenv("GOPATH"), "src/github.com/kvlinden/go/hello-plugin/templates/*"))
	// A simple index page to be served up by goCMS at /.
	router.GET("/", helloPlugin)
	// This "public" path needs to have /api so that it matches what goCMS does. TODO: Is this not a bad dependency?
	router.GET("/api/public", helloPublic)
	// Use a gin group to bundle pages that require authorization. This /api base path must, again, match goCMS.
	authGroup := router.Group("/api")
	// Add a sample middleware function that could be used to authenticate or log.
	authGroup.Use(helloMiddleware)
	authGroup.GET("/auth", helloAuth)

	router.Run(fmt.Sprintf("localhost: %v", *port))
}

func helloPlugin(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{"title": "Hello Plugin"})
}

func helloPublic(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello, public page!",
	})
}

func helloAuth(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello, auth page!",
	})
}

func helloMiddleware(c *gin.Context) {
	log.Println("Hello, middleware!")
	c.Next()
}