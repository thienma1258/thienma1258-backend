package services

import (
	"dongpham/model"
	"dongpham/repository"
)

type PostServices struct {
	Repo *repository.PostRepository
}

func (gs *PostServices) GetAllPostIDs() ([]int, error) {
	ats, err := gs.Repo.GetAllPostIDs()
	return ats, err
}

func (gs *PostServices) Create(post model.Post) error {
	return gs.Repo.CreateNewPost(&post)
}

func NewPostServices(repo *repository.PostRepository) *PostServices {
	return &PostServices{Repo: repo}
}
