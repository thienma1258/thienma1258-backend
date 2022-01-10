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

type QueryPost struct {
	Published *bool
}

var POST_COLUMNS = []string{"id", "title", "slug", "image", "body", "published", "created_at", "updated_at"}

func NewPostRepository(db *sql.DB) (*PostRepository, error) {
	return &PostRepository{db: db}, nil
}

func (up *PostRepository) GetAllPostIDs(query QueryPost) ([]int, error) {
	//up.db.Exec()
	sqlBuilder := psql.Select("id").From("posts")
	if query.Published != nil {
		sqlBuilder = sqlBuilder.Where(sq.Eq{"published": query.Published})
	}

	sqlQuery, _, err := sqlBuilder.ToSql()
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

func (up *PostRepository) GetPostByIDs(ids []int) ([]*model.Post, error) {
	//up.db.Exec()\
	//sqlQuery, args, err := sq.Select("id")

	var result []*model.Post
	sqlBuilder := psql.Select(POST_COLUMNS...).From("posts")
	var args []interface{}
	if len(ids) > 0 {
		for _, id := range ids {
			args = append(args, id)
		}

		sqlBuilder = sqlBuilder.Where("id in (?"+strings.Repeat(",?", len(ids)-1)+")", args...)
	}
	sqlQuery, _, err := sqlBuilder.ToSql()
	if err != nil {
		return result, err
	}

	rows, err := up.db.Query(sqlQuery, args...)
	if err != nil {
		return result, err
	}

	for rows.Next() {
		post := &model.Post{}
		err = rows.Scan(&post.ID, &post.Title, &post.Slug, &post.Image, &post.Body, &post.Published, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return result, nil
		}
		result = append(result, post)
	}

	return result, nil
}

func (up *PostRepository) CreateNewPost(newPost *model.Post) error {
	//up.db.Exec()

	newPost.CreatedAt = utils.String(utils.GetCurrentISOTime())
	newPost.UpdatedAt = utils.String(utils.GetCurrentISOTime())
	sqlQuery, args, err := psql.
		Insert("posts").Columns("title", "slug", "image", "body", "published", "created_at", "updated_at").
		Values(newPost.Title, newPost.Slug, newPost.Image, newPost.Body, newPost.Published, newPost.CreatedAt, newPost.UpdatedAt).
		ToSql()
	_, err = up.db.Exec(sqlQuery, args...)
	return err
}

func (up *PostRepository) UpdatePost(updatePost *model.Post) error {
	sqlQuery, args, err := psql.
		Update("posts").
		Set("title", updatePost.Title).
		Set("slug", updatePost.Slug).
		Set("image", updatePost.Image).
		Set("body", updatePost.Body).
		Set("published", updatePost.Published).
		Set("updatedAt", utils.GetCurrentTimeStamp()).
		ToSql()
	_, err = up.db.Exec(sqlQuery, args)
	return err
}
