package services

import (
	"dongpham/model"
	"dongpham/repository"
)

type PostServices struct {
	Repo *repository.PostRepository
}

func (gs *PostServices) GetAllPostIDs(published *bool) ([]int, error) {
	ats, err := gs.Repo.GetAllPostIDs(repository.QueryPost{Published: published})
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


func (gs *PostServices) Create(post model.Post) error {
	return gs.Repo.CreateNewPost(&post)
}

func (gs *PostServices) Update(post model.Post) error {
	return gs.Repo.UpdatePost(&post)
}

func NewPostServices(repo *repository.PostRepository) *PostServices {
	return &PostServices{Repo: repo}
}
