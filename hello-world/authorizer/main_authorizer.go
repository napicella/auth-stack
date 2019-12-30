package main

import (
	"com.napicella/hello-world/utils"
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"strings"
)


func handleRequest(ctx context.Context, event events.APIGatewayCustomAuthorizerRequestTypeRequest) (
	events.APIGatewayCustomAuthorizerResponse, error) {

	fmt.Println("==== Authorizer hello world ====")

	tmp := strings.Split(event.MethodArn, ":")
	apiGatewayArnTmp := strings.Split(tmp[5], "/")
	awsAccountID := tmp[4]
	resp := NewAuthorizerResponse("", awsAccountID)
	resp.Region = tmp[3]
	resp.APIID = apiGatewayArnTmp[0]
	resp.Stage = apiGatewayArnTmp[1]

	resp.PrincipalID = "unknown"

	if strings.HasPrefix(event.Path, "/s") {
		fmt.Println("path starts with /s - everybody allowed here")
		fmt.Println(event.Path)
		// TODO: It was put in place 'cause of this, but maybe this is not necessary:
		// https://stackoverflow.com/questions/50331588/aws-api-gateway-custom-authorizer-strange-showing-error
		resp.addMethod(Allow, All, "/s/*")
		return resp.APIGatewayCustomAuthorizerResponse, nil
	}

	fmt.Println("path does not start with /s - expecting to find valid jwt token in the Cookie:token=...")

	token := getCookie(event, "token")
	fmt.Printf("token: %s\n", token)
	if token == "" {
		fmt.Println("no token found in the request")
		resp.DenyAllMethods()
		return resp.APIGatewayCustomAuthorizerResponse, errors.New("no token in the request")
	}
	secret, e := utils.GetSecret()
	if e != nil {
		fmt.Println("unable to get the secret")
		fmt.Println(e)
		resp.DenyAllMethods()
		return resp.APIGatewayCustomAuthorizerResponse, errors.New("unable to parse token")
	}
	c, e := utils.VerifyToken(token, secret)
	if e != nil {
		fmt.Println("failed to verify token")
		fmt.Println(e)
		resp.DenyAllMethods()
		return resp.APIGatewayCustomAuthorizerResponse, errors.New("unable to verify token")
	}
	fmt.Println("everything is fine, allowing the request")
	resp.PrincipalID = c.Username
	resp.addMethod(Allow, All, "/api/*")
	return resp.APIGatewayCustomAuthorizerResponse, nil
}

func getCookie(event events.APIGatewayCustomAuthorizerRequestTypeRequest, cookieName string) string {
	cookies := event.MultiValueHeaders["Cookie"]
	fmt.Println(cookies)
	for _, cookie := range cookies {
		fmt.Println(cookie)
		parts := strings.Split(cookie, "=")
		if parts[0] == cookieName {
			return parts[1]
		}
	}
	return ""
}

func main() {
	lambda.Start(handleRequest)
}