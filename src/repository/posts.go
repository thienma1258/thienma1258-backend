package repository

import (
	"database/sql"
	"dongpham/model"
	"dongpham/utils"
	sq "github.com/Masterminds/squirrel"
	"strings"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)


type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) (*PostRepository, error) {
	return &PostRepository{db: db}, nil
}

func (up *PostRepository) GetAllPostIDs() ([]int, error) {
	//up.db.Exec()
	sqlQuery, _, err := sq.Select("id").From("posts").ToSql()
	if err != nil {
		return []int{}, err
	}

	rows, err := up.db.Query(sqlQuery)
	if err != nil {
		return []int{}, err
	}

	ids := ScanRowIDs(*rows)

	return ids, nil
}

func (up *PostRepository) GetPostByIDs(ids []string) ([]*model.Post, error) {
	//up.db.Exec()\
	//sqlQuery, args, err := sq.Select("id")

	var result []*model.Post
	sqlBuilder := sq.Select("*").From("posts")
	if len(ids) > 0 {
		sqlBuilder.Where("id in ("+strings.TrimSuffix(strings.Repeat("?,", len(ids))+")", ","), ids)
	}
	sqlQuery, args, err := sqlBuilder.ToSql()
	if err != nil {
		return result, err
	}

	rows, err := up.db.Query(sqlQuery, args)
	if err != nil {
		return result, err
	}

	for rows.Next() {
		post := &model.Post{}
		err = rows.Scan(&post.ID, &post.Title, &post.Image, &post.UpdatedAt, &post.CreatedAt)
		if err != nil {
			return result, nil
		}
		result = append(result, post)
	}

	return result, nil
}

func (up *PostRepository) CreateNewPost(newPost *model.Post) error {
	//up.db.Exec()

	newPost.CreatedAt = utils.GetCurrentTimeStamp()
	newPost.UpdatedAt = utils.GetCurrentTimeStamp()
	sqlQuery, args, err := psql.
		Insert("posts").Columns("title", "slug", "image", "body", "published", "created_at", "updated_at").
		Values(newPost.Title, newPost.Slug, newPost.Image, newPost.Body, newPost.Published, newPost.CreatedAt, newPost.UpdatedAt).
		ToSql()
	_, err = up.db.Exec(sqlQuery, args...)
	return err
}

func (up *PostRepository) UpdatePost(updatePost *model.Post) error {
	sqlQuery, args, err := sq.
		Update("posts").
		Set("title", updatePost.Title).
		Set("slug", updatePost.Slug).
		Set("image", updatePost.Image).
		Set("body", updatePost.Body).
		Set("published", updatePost.Published).
		Set("updatedAt", updatePost.UpdatedAt).
		ToSql()
	_, err = up.db.Exec(sqlQuery, args)
	return err
}
