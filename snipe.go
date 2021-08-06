package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

func getToken(email string, password string) string {
	type ResJSON struct {
		ClientToken string `json:"clientToken"`
		AccessToken string `json:"accessToken"`
	}

	url := "https://authserver.mojang.com/authenticate"
	var jsonStr = []byte(`{"username":"` + email + `", "password": "` + password + `"}`)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	res, _ := http.DefaultClient.Do(req)
	if res.StatusCode == 200 {
		// Logged in!
		success("Logged in!")
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			error("Error reading access token!")
			os.Exit(-1)
		}
		var rj ResJSON
		json.Unmarshal([]byte(string(body)), &rj)
		return rj.AccessToken
	} else {
		fmt.Println(res)
	}

	return ""
}

func getTokenMS(email string, password string) string {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			Renegotiation:      tls.RenegotiateOnceAsClient,
			InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
	}
	res, err := http.Get("https://login.live.com/oauth20_authorize.srf?client_id=000000004C12AE6F&redirect_uri=https://login.live.com/oauth20_desktop.srf&scope=service::user.auth.xboxlive.com::MBI_SSL&display=touch&response_type=token&locale=en")
	if err != nil {
		error("Unable to login, try again later!")
		os.Exit(-1)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		error("Unable to login, try again later!")
		os.Exit(-1)
	}
	sb := string(body)
	sFTTag := ""
	regex, err := regexp.Compile("value=\"(.+?)\"")
	if err != nil {
		error("Error parsing login data!")
		os.Exit(-1)
	}
	sFTTag = regex.FindStringSubmatch(sb)[0]
	urlPost := ""
	regex2, err := regexp.Compile("urlPost:'(.+?)'")
	if err != nil {
		error("Error parsing login data!")
		os.Exit(-1)
	}
	urlPost = regex2.FindStringSubmatch(sb)[0]

	sFTTag = strings.Split(sFTTag, "=")[1]
	sFTTag = strings.ReplaceAll(sFTTag, "\"", "")
	urlPost = strings.Join(strings.Split(urlPost, ":")[1:], ":")
	urlPost = strings.ReplaceAll(urlPost, "'", "")

	//fmt.Println("urlPost: " + urlPost)
	req, _ := http.NewRequest("POST", urlPost, bytes.NewBuffer([]byte("login="+url.QueryEscape(email)+"&loginfmt="+url.QueryEscape(email)+"&passwd="+url.QueryEscape(password)+"&PPFT="+url.QueryEscape(sFTTag)+"&type=11&LoginOptions=3&i13=0&ps=2&PPSX=Pass")))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err = client.Do(req)
	if err != nil {
		fmt.Println(err)
		error("Error logging you in!")
		os.Exit(-1)
	}
	finalURL := res.Request.URL.String()
	//fmt.Println(finalURL)
	if strings.Contains(finalURL, "accessToken") && (strings.Contains(finalURL, urlPost) || strings.Contains(urlPost, finalURL)) {
		error("Error fetching token from login URL!")
		os.Exit(-1)
	}
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		error("Error parsing response body!")
		os.Exit(-1)
	}
	sb = string(body)
	if strings.Contains(sb, "Sign in to") {
		error("Invalid credentials! Make sure to enter your correct email and password!")
		enterManually := input("Would you like to complete the setup manually? (Y/N) ")
		if strings.ToLower(enterManually) == "y" {
			info("Follow the URL in the next line, enter your password and press SIGN IN. Then copy the new URL in the address bar and paste it here!")
			info(urlPost)
			finalURL = input("New URL: ")
		} else {
			os.Exit(-1)
		}
	}
	if strings.Contains(sb, "Help us protect your account") {
		error("2FA present, unable to login!")
		os.Exit(-1)
	}
	type MSLoginResponse struct {
		AccessToken string
		TokenType   string
		ExpiresIn   string
	}
	var authData MSLoginResponse
	queryParams := strings.Split(finalURL, "#")[1]
	paramPairs := strings.Split(queryParams, "&")
	for _, param := range paramPairs {
		if strings.Contains(param, "access_token") {
			authData.AccessToken, _ = url.QueryUnescape(strings.Split(param, "=")[1])
		} else if strings.Contains(param, "token_type") {
			authData.TokenType, _ = url.QueryUnescape(strings.Split(param, "=")[1])
		} else if strings.Contains(param, "expires_in") {
			authData.ExpiresIn, _ = url.QueryUnescape(strings.Split(param, "=")[1])
		}
	}
	req, _ = http.NewRequest("POST", "https://user.auth.xboxlive.com/user/authenticate", bytes.NewBuffer([]byte(`{"Properties":{"AuthMethod":"RPS","SiteName":"user.auth.xboxlive.com","RpsTicket":"`+authData.AccessToken+`"},"RelyingParty":"http://auth.xboxlive.com","TokenType":"JWT"}`)))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	res, err = client.Do(req)
	if err != nil {
		error("Error signing into Xbox Live!")
		error(err.Error())
		//fmt.Println(res)
		os.Exit(-1)
	}
	body, err = ioutil.ReadAll(res.Body)
	type XboxResponse struct {
		Token         string `json:"Token"`
		DisplayClaims struct {
			Xui []struct {
				Uhs string `json:"uhs"`
			} `json:"xui"`
		} `json:"DisplayClaims"`
	}
	var xbr XboxResponse
	json.Unmarshal([]byte(string(body)), &xbr)

	req, _ = http.NewRequest("POST", "https://xsts.auth.xboxlive.com/xsts/authorize", bytes.NewBuffer([]byte(`{"Properties":{"SandboxId":"RETAIL","UserTokens":["`+xbr.Token+`"]},"RelyingParty":"rp://api.minecraftservices.com/","TokenType":"JWT"}`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	res, err = client.Do(req)
	if err != nil || res.StatusCode == 401 {
		error("Error getting response from Xbox Live!")
		os.Exit(-1)
	}
	body, err = ioutil.ReadAll(res.Body)
	type NewXboxTokenResponse struct {
		Token string `json:"Token"`
	}
	var nxbr NewXboxTokenResponse
	json.Unmarshal([]byte(string(body)), &nxbr)

	req, _ = http.NewRequest("POST", "https://api.minecraftservices.com/authentication/login_with_xbox", bytes.NewBuffer([]byte(`{"identityToken":"XBL3.0 x=`+xbr.DisplayClaims.Xui[0].Uhs+`;`+nxbr.Token+`","ensureLegacyEnabled":true}`)))
	req.Header.Set("Content-Type", "application/json")
	res, err = client.Do(req)
	if err != nil {
		error("Error fetching bearer token!")
		os.Exit(-1)
	}
	type MinecraftTokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
	}
	body, err = ioutil.ReadAll(res.Body)
	var mtr MinecraftTokenResponse
	json.Unmarshal([]byte(string(body)), &mtr)
	return mtr.AccessToken
}

func validateToken(token string) bool {
	url := "https://api.minecraftservices.com/minecraft/profile"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	res, _ := http.DefaultClient.Do(req)
	return res.StatusCode != 401
}

func changeSkin(token string) {
	var jsonStr = []byte(`{"variant": "classic", "url": "https://raw.githubusercontent.com/HumanDuck23/goSniper/master/goSniper.png"}`)
	url := "https://api.minecraftservices.com/minecraft/profile/skins"
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	_, err := http.DefaultClient.Do(req)
	if err != nil {
		error("Error changing skin!")
	}
}

func getDropTime(username string) int {
	url := "http://api.coolkidmacho.com/droptime/" + username
	res, err := http.Get(url)
	if err != nil {
		return -1
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return -1
	}
	var data DropTime
	json.Unmarshal([]byte(string(body)), &data)
	if data.UNIX == 0 { // If the API returns 0 as the droptime, fail
		return -1
	}
	return data.UNIX
}

func snipe(config SnipeConfig, group *sync.WaitGroup) {
	info(fmt.Sprintf("Sniping %s at %s / in %s ms!", magenta(config.USERNAME, true), magenta(strconv.Itoa(config.DROPTIME), true), magenta(strconv.Itoa((config.DROPTIME-int(time.Now().Unix()))*1000), true)))
	ping := getPing()
	info("Ping: " + strconv.Itoa(int(ping)))
	timeAt := int64((float64(config.DROPTIME) - (float64(config.OFFSET) / 1000) - (ping / 1000)) * 1e9)
	url_ := "https://api.minecraftservices.com/minecraft/profile/name/" + config.USERNAME
	req, _ := http.NewRequest("PUT", url_, nil)
	req.Header.Set("Authorization", "Bearer "+config.TOKEN)
	for i := 0; i < 3; i++ {
		group.Add(1)
		tmpOffset := int64(i * 90 * 1e6)
		go sendRequest(config.TOKEN, timeAt+tmpOffset, req, group)
	}
	defer group.Done()
}

func sendRequest(token string, timeAt int64, r *http.Request, group *sync.WaitGroup) {
	for time.Now().UnixNano() < timeAt {
		time.Sleep(time.Microsecond)
	}
	timeInfo("Attempting to snipe now!")
	res, _ := http.DefaultClient.Do(r)
	if res.StatusCode == 200 {
		success("Name sniped!")
		fmt.Println(time.Now())
		changeSkin(token)
	}
	log(fmt.Sprint(time.Now()) + " ==== " + fmt.Sprint(res))
	defer group.Done()
}

func getPing() float64 {
	start := time.Now().UnixNano()
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.minecraftservices.com/minecraft/profile/name/abcdef", nil)
	_, _ = client.Do(req)
	end := time.Now().UnixNano()
	return float64(end-start) / 1e6
}
