package services

import (
	"dongpham/model"
	"dongpham/repository"
	"dongpham/utils"
	"github.com/gosimple/slug"
)

type PostServices struct {
	Repo *repository.PostRepository
}

const DEFAULT_AUTHOR = "thienma1258"

func (gs *PostServices) GetAllPostIDs(published *bool,orderDesc *bool) ([]int, error) {
	ats, err := gs.Repo.GetAllPostIDs(repository.QueryPost{
		Published: published,
		OrderDESC: orderDesc,
	})
	return ats, err
}

func (gs *PostServices) GetPostByIDs(ids []int, _fields []string) (map[int]*model.Post, error) {
	ats, err := gs.Repo.GetPostByIDs(ids)

	if err != nil {
		return nil, err
	}
	result := map[int]*model.Post{}
	for key, val := range ats {
		result[key] = val
	}
	return result, err
}

func (gs *PostServices) Create(post model.Post) (int,error) {
	beforeInsertOfUpdate(&post)
	return gs.Repo.CreateNewPost(&post)
}

func (gs *PostServices) Update(post model.Post) error {
	beforeInsertOfUpdate(&post)
	return gs.Repo.UpdatePost(&post)
}

func NewPostServices(repo *repository.PostRepository) *PostServices {
	return &PostServices{Repo: repo}
}

func beforeInsertOfUpdate(post *model.Post) {
	if post.Slug == nil || len(*post.Slug) == 0 {
		post.Slug = utils.String(slug.Make(utils.StringValue(post.Title)))
	}

	if post.SocialDescription == nil || len(*post.SocialDescription) == 0 {
		post.SocialDescription = post.Description
	}

	if post.SocialTitle == nil || len(*post.SocialTitle) == 0 {
		post.SocialTitle = post.Title
	}

	if post.SocialImage == nil || len(*post.SocialImage) == 0 {
		post.Image = post.Image
	}

	if post.Author == nil || len(*post.Author) == 0 {
		post.Author = utils.String(DEFAULT_AUTHOR)
	}

}
