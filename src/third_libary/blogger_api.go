package third_libary

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"log"
	"net/http"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const API_SEARCH_END_POINT_FORMAT = "https://www.googleapis.com/blogger/v3/blogs/3213900/posts/search?q=$%s&key=%s"

type BloggerResponse struct {
	Kind  string        `json:"kind"`
	Items []BloggerData `json:"items"`
}

type BloggerData struct {
	Kind     string `json:"kind"`
	Id       string `json:"id"`
	Updated  string `json:"updated"`
	Url      string `json:"url"`
	SelfLink string `json:"selfLink"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Author   struct {
		DisplayName string `json:"displayName"`
		Url         string `json:"url"`
	}
}

func GetSearchNewPost(q string, accessKey string) ([]BloggerData, error) {
	resp, err := http.Get(fmt.Sprintf(API_SEARCH_END_POINT_FORMAT, q, accessKey))
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var result BloggerResponse

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result.Items, nil
}
