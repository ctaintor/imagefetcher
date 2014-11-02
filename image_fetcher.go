// go run boss.go --consumerkey <key> --consumersecret <secret> --appname <appname>
package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
  "github.com/go-martini/martini"

	"github.com/mrjones/oauth"
)

type yahooResponse struct {
	BossResponse bossResponse `json:"bossresponse"`
}

type bossResponse struct {
	Images bossImages `json:"images"`
}

type bossImages struct {
	Results []bossImageResult `json:"results"`
}

type bossImageResult struct {
	Url    string `json:"url"`
	Format string `json:"format"`
}

func getImageUrls(unparsedResponse []byte) []string {
	parsedResponse := yahooResponse{}
	json.NewDecoder(strings.NewReader(string(unparsedResponse))).Decode(&parsedResponse)

	images := make([]string, 0, 35)
	for i := range parsedResponse.BossResponse.Images.Results {
		imageResult := parsedResponse.BossResponse.Images.Results[i]
		if imageResult.Format == "jpeg" || imageResult.Format == "jpg" ||
			imageResult.Format == "gif" || imageResult.Format == "png" {
			images = append(images, imageResult.Url)
		}
	}

	return images
}

func main() {
	c := oauth.NewConsumer(
		os.Getenv("CONSUMER_KEY"),
		os.Getenv("CONSUMER_SECRET"),
		oauth.ServiceProvider{
			RequestTokenUrl:   "https://api.login.yahoo.com/oauth/v2/get_request_token",
			AuthorizeTokenUrl: "https://api.login.yahoo.com/oauth/v2/request_auth",
			AccessTokenUrl:    "https://api.login.yahoo.com/oauth/v2/get_token",
		})

	accessToken := &oauth.AccessToken{}
	accessToken.Token = os.Getenv("ACCESS_TOKEN")
	accessToken.Secret = os.Getenv("TOKEN_SECRET")
	accessToken.AdditionalData = map[string]string{"oauth_expires_in": "3600", "oauth_session_handle": os.Getenv("SESSION_HANDLE"), "oauth_authorization_expires_in": "732555938", "xoauth_yahoo_guid": os.Getenv("YAHOO_GUID")}

	m := martini.Classic()

	//curl -X POST  -k -u token: "http://localhost:3000/get_image.json" -d '{"urls": ["http://talks.golang.org/2013/advconc/gopherswim.jpg", "http://www.unixstickers.com/image/cache/data/stickers/golang/golang.sh-600x600.png"], "height_px": 400}' > test.png
	m.Get("/v1/imageUrls/:word", func(params martini.Params) []byte {

		response, err := c.Get(
			"https://yboss.yahooapis.com/ysearch/images",
			map[string]string{"q": params["word"], "sites": "", "format": "json", "dimensions": "medium"},
			accessToken)
		if err != nil {
			log.Fatal(err)
		}
		defer response.Body.Close()

		bits, err := ioutil.ReadAll(response.Body)

		imageUrls := getImageUrls(bits)

		jsonString, _ := json.Marshal(imageUrls)

		return jsonString
	})

	m.Run()
}
