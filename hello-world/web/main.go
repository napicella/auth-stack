package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/apex/gateway"
	"github.com/gin-gonic/gin"
	"github.com/mdp/qrterminal/v3"
	"github.com/xlzd/gotp"
)

func welcomeHandler(c *gin.Context) {
	c.String(http.StatusOK, "Hello World from Go")
}

func rootHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{ "message" : "hello"})
}

func routerEngine() (*gin.Engine, error) {
	// set server mode
	gin.SetMode(gin.DebugMode)

	r := gin.New()

	// Global middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())


	signinHandler, e := Signin()
	if e != nil {
		return nil, e
	}
	r.GET("/", rootHandler)

	// "/s" apis are not to be protected by authentication
	r.POST("/s/signin", signinHandler)
	r.POST("/s/user/register", RegisterUser)
	r.GET("/s/welcome-2", welcomeHandler)

	r.GET("/s/ping", pingHandler)
	r.POST("/api/token/refresh", Refresh)
	r.GET("/api/totp/gen", GenerateTotp)
	r.POST("/api/totp/register", RegisterTotp)

	return r, nil
}

func pingHandler(context *gin.Context) {
	r := context.Request
	// example retrieving values from the api gateway proxy request context.
	requestContext, ok := gateway.RequestContext(r.Context())
	if !ok || requestContext.Authorizer["name"] == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message" : "name context value not found"})
		return
	}
	userID := requestContext.Authorizer["name"].(string)
	fmt.Println(userID)
	context.JSON(http.StatusOK, gin.H{"message" : userID})
}

type registerTotp struct {
	Totp string `json:"totp"`
}

func RegisterTotp(c *gin.Context) {
	var totpRequest registerTotp
	e := c.ShouldBindJSON(&totpRequest)
	if e != nil {
		fmt.Println("error while registering the totp")
		fmt.Println(e)
		c.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
		return
	}
	totp := gotp.NewDefaultTOTP("JBSWY3DPEHPK3PXP")
	isValid := totp.Verify(totpRequest.Totp, int(time.Now().Unix()))
	if isValid {
		fmt.Println("totp is valid")
		c.JSON(http.StatusOK, gin.H{"message" : "valid totp"})
		return
	}
	fmt.Println("error while registering the totp")
	fmt.Println("totp is invalid")
	c.JSON(http.StatusBadRequest, gin.H{"message" : "invalid totp"})}

func GenerateTotp(context *gin.Context) {
	var buff bytes.Buffer
	qrterminal.Generate("otpauth://totp/Example:alice@google.com?secret=JBSWY3DPEHPK3PXP&issuer=Example",
		qrterminal.L,
		&buff)
	context.Writer.WriteString(buff.String())
	context.Writer.WriteHeader(http.StatusOK)
}

func main() {
	addr := ":" + os.Getenv("PORT")
	router, e := routerEngine()
	if e != nil {
		log.Fatal(e)
	}
	log.Fatal(gateway.ListenAndServe(addr, router))
}