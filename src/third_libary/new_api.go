package third_libary

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

//GET https://newsapi.org/v2/everything?q=bitcoin&apiKey=2112942a1c794515af9ab1d9aaff3263
const API_NEWS_SEARCH_END_POINT_FORMAT = "https://newsapi.org/v2/everything?q=%s&apiKey=%s"

type NewResponseAPI struct {
	Status      string `json:"status"`
	TotalResult int    `json:"totalResult"`

	Items []Article `json:"articles"`
}

type Article struct {
	Author      string `json:"author"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"urlToImage"`
	PublishedAt string `json:"publishedAt"`
	Content     string `json:"content"`
	URL         string `json:"url"`
}

func GetSearchNewPostFromNewAPI(q string, accessKey string, from string) ([]Article, error) {
	resp, err := http.Get(fmt.Sprintf(API_NEWS_SEARCH_END_POINT_FORMAT, q, accessKey))
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var result NewResponseAPI

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result.Items, nil
}
