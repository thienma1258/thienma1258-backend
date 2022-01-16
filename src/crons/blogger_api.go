package crons

import (
	"dongpham/config"
	"dongpham/model"
	"dongpham/services"
	"dongpham/third_libary"
	"dongpham/utils"
)

func GetDataFromBloggerAPI() error {
	if len(config.BloggerAPIKey) > 0 {
		data, err := third_libary.GetSearchNewPost("api", config.BloggerAPIKey)
		if err != nil {
			return err
		}
		if len(data) == 0 {
			return nil
		}

		var insertNewPost []*model.Post

		for _, item := range data {
			insertNewPost = append(insertNewPost, ParseToPostModal(&item))
		}
		err = services.GetPostService().CreateMultiplePost(insertNewPost)
		return err
	}
	return nil
}

func ParseToPostModal(data *third_libary.BloggerData) *model.Post {
	return &model.Post{
		Meta: map[string]interface{}{
			"url":      data.Url,
			"selfLink": data.SelfLink,
			"id":       data.Id,
			"kind":     data.Kind,
		},
		Author:    utils.String(data.Author.DisplayName),
		Body:      utils.String(data.Content),
		Published: utils.Bool(false),
		Title:     utils.String(data.Title),
		Image:     utils.String("https://pipedream.com/s.v0/app_OkrhoY/logo/orig"),
	}
}
