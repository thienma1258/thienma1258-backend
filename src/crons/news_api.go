package crons

import (
	"dongpham/config"
	"dongpham/model"
	"dongpham/services"
	"dongpham/third_libary"
	"dongpham/utils"
)

func GetDataFromNew() error {
	if len(config.NewApiKey) > 0 {
		data, err := third_libary.GetSearchNewPostFromNewAPI("sql", config.NewApiKey, "2022-01-01")
		if err != nil {
			return err
		}
		if len(data) == 0 {
			return nil
		}

		var insertNewPost []*model.Post

		for _, item := range data {
			insertNewPost = append(insertNewPost, ParseFromNewArticleToPostModal(&item))
		}
		err = services.GetPostService().CreateMultiplePost(insertNewPost)
		return err
	}
	return nil
}

func ParseFromNewArticleToPostModal(data *third_libary.Article) *model.Post {
	return &model.Post{
		Meta: map[string]interface{}{
			"url": data.URL,
		},
		Author:      utils.String(data.Author),
		Body:        utils.String(data.Content),
		Published:   utils.Bool(false),
		Title:       utils.String(data.Title),
		Image:       utils.String(data.Image),
		Description: utils.String(data.Description),
	}
}
