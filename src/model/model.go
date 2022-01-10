package model

import (
	"io"
	"net/url"
	"time"
)

type User struct {
	Id             string
	UserName       string
	Profile        string
	FullName       string
	Introduction   string
	DateOfBirth    string
	LastUpdateTime time.Time
}

type GalleryImage struct {
	Id               string
	Description      string
	OriginalImageUrl string
	ResizeImageUrl   string
	Folder           string
	Order            int
	OwnerId          string
	LastUpdateTime   time.Time
}

type Article struct {
	Id             string
	BodyHtml       string
	OwnerId        string
	PictureUrl     string
	Description    string
	Priority       int
	Order          int
	Attachment     string
	LastUpdateTime time.Time
}

type Carousel struct {
	Id             string
	PictureUrl     string
	OwnerId        string
	OrderId        int
	LinkUrl        string
	LastUpdateTime time.Time
}

type ApiRequest struct {
	Body      io.Reader
	AuthToken string
	Query     map[string]string
	URLQuery  url.Values
}

type ApiResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}
