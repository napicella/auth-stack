package main

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var cred = `{ "password" : "password1", "username" : "user1" }`
const url = "https://DOMAIN_NAME_HERE.execute-api.us-west-2.amazonaws.com/Prod/"

func getUrl(path string) string {
	return url + path
}

func main()  {
	createUser()

	req, err := http.NewRequest(
		"POST",
		getUrl("/s/signin"),
		strings.NewReader(cred))

	if err != nil {
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	var tokenCookie *http.Cookie
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "token" {
			tokenCookie = cookie
			break
		}
	}

	req, err = http.NewRequest(
		"GET",
		getUrl("/api/ping"), nil)
	if err != nil {
		return
	}

	req.AddCookie(tokenCookie)
	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = errors.New("resp.StatusCode: " + strconv.Itoa(resp.StatusCode))
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))


	req, err = http.NewRequest(
		"GET",
		getUrl("/api/totp/gen"), nil)
	if err != nil {
		return
	}

	req.AddCookie(tokenCookie)
	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = errors.New("resp.StatusCode: " + strconv.Itoa(resp.StatusCode))
		return
	}

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))


	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter totp token: ")
	text, _ := reader.ReadString('\n')
	fmt.Println(strings.ReplaceAll(strings.TrimSpace(text), "\n", ""))

	sendTotp(tokenCookie, strings.ReplaceAll(strings.TrimSpace(text), "\n", ""))
}

func sendTotp(tokenCookie *http.Cookie, totp string)  {
	fmt.Println(fmt.Sprintf(`{ "totp" : "%s" }`, totp))
	req, err := http.NewRequest(
		http.MethodPost,
		getUrl("/api/totp/register"),
		strings.NewReader(fmt.Sprintf(`{ "totp" : "%s" }`, totp)))
	if err != nil {
		return
	}

	req.AddCookie(tokenCookie)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func createUser()  {
	req, err := http.NewRequest(
		"POST",
		getUrl("/s/user/register"),
		strings.NewReader(`{"username" : "user1", "password" : "password1" }`))

	if err != nil {
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = errors.New("resp.StatusCode: " + strconv.Itoa(resp.StatusCode))
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
