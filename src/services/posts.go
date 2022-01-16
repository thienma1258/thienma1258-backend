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

var postService *PostServices

func init() {
	postService = &PostServices{Repo: repository.PostRepo}
}

func GetPostService() *PostServices {
	return postService
}

const DEFAULT_AUTHOR = "thienma1258"
const DEFAULT_USER = "system"

func (gs *PostServices) GetAllPostIDs(published *bool, orderDesc *bool) ([]int, error) {
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
	for _, val := range ats {
		result[utils.IntVal(val.ID)] = val
	}
	return result, err
}

func (gs *PostServices) Create(post model.Post) (int, error) {
	beforeInsertOfUpdate(&post)
	return gs.Repo.CreateNewPost(&post)
}

func (gs *PostServices) Update(post model.Post) error {
	beforeInsertOfUpdate(&post)
	return gs.Repo.UpdatePost(&post)
}

func (gs *PostServices) Delete(id int) error {
	return gs.Repo.Delete(id)
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
		post.SocialImage = post.Image
	}

	if post.Author == nil || len(*post.Author) == 0 {
		post.Author = utils.String(DEFAULT_AUTHOR)
	}

	if post.UserID == nil || len(*post.UserID) == 0 {
		post.UserID = utils.String(DEFAULT_USER)
	}

}

func (gs *PostServices) CreateMultiplePost(posts []*model.Post) error {
	mapName := map[string]struct{}{}
	var createPosts []*model.Post
	for _, post := range posts {
		beforeInsertOfUpdate(post)
		if _, ok := mapName[*post.Slug]; !ok {
			mapName[*post.Slug] = struct{}{}
			createPosts = append(createPosts, post)
		} else {
			//return fmt.Errorf("name is duplicate %s", *post.Slug)
		}
	}
	return gs.Repo.CreateMultipleNewPost(createPosts)
}
