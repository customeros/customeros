package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func main() {

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "OPTIONS"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))

	router.GET("/customer-os/registered-apps", getCustomerOsRegisteredApps)
	router.Run("localhost:8070")
}

type RegisteredApp struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type AuthenticatedResponse struct {
	CONTENT []byte          `json:"content"`
	APPS    []RegisteredApp `json:"apps"`
}

func getCustomerOsRegisteredApps(context *gin.Context) {
	//content, err := getUserInfo(context.Request.FormValue("state"), context.Request.FormValue("code"))
	//if err != nil {
	//	context.Redirect(http.StatusTemporaryRedirect, "/")
	//	return
	//}
	var apps = []RegisteredApp{
		{ID: "1", Name: "Oasis", URL: "http://oasis.openline.local"},
		{ID: "2", Name: "Contacts", URL: "http://contacts.openline.local"},
	}
	var response AuthenticatedResponse

	response.APPS = apps
	response.CONTENT = []byte("content")
	context.IndentedJSON(http.StatusOK, response)
}
