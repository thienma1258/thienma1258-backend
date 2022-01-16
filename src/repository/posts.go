package repository

import (
	"database/sql"
	"dongpham/internal_errors"
	"dongpham/model"
	"dongpham/utils"
	"encoding/json"
	sq "github.com/Masterminds/squirrel"
	"log"
	"strings"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type PostRepository struct {
	db *sql.DB
}

type QueryPost struct {
	Published *bool
	OrderDESC *bool
}

var POST_COLUMNS = []string{"id", "title", "slug", "description", "image", "body", "published", "created_at", "updated_at", "social_title", "social_description", "social_image", "author", "meta"}

func NewPostRepository(db *sql.DB) (*PostRepository, error) {
	return &PostRepository{db: db}, nil
}

func (up *PostRepository) GetAllPostIDs(query QueryPost) ([]int, error) {
	//up.db.Exec()
	sqlBuilder := psql.Select("id").From("posts")
	if query.Published != nil {
		sqlBuilder = sqlBuilder.Where(sq.Eq{"published": query.Published})
	}
	if query.OrderDESC == nil {
		sqlBuilder = sqlBuilder.OrderBy("updated_at DESC")
	} else {
		orderBy := ""
		if *query.OrderDESC {
			orderBy = "DESC"
		} else {
			orderBy = "ASC"
		}
		sqlBuilder = sqlBuilder.OrderBy("updated_at " + orderBy)
	}
	sqlBuilder = sqlBuilder.Where(sq.Eq{"deleted": false})
	sqlQuery, aggs, err := sqlBuilder.ToSql()
	if err != nil {
		return []int{}, err
	}

	rows, err := up.db.Query(sqlQuery, aggs...)
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
		var meta sql.NullString
		err = rows.Scan(&post.ID, &post.Title, &post.Slug, &post.Description, &post.Image, &post.Body, &post.Published, &post.CreatedAt, &post.UpdatedAt,
			&post.SocialTitle, &post.SocialDescription, &post.Image, &post.Author, &meta)
		if meta.Valid {

			err = json.Unmarshal([]byte(meta.String), &post.Meta)
			if err != nil {
				return result, err
			}
		}

		if err != nil {
			return result, nil
		}
		result = append(result, post)
	}

	return result, nil
}

func (up *PostRepository) CreateNewPost(newPost *model.Post) (int, error) {
	//up.db.Exec()

	newPost.CreatedAt = utils.String(utils.GetCurrentISOTime())
	newPost.UpdatedAt = utils.String(utils.GetCurrentISOTime())
	meta, err := json.Marshal(newPost.Meta)
	if err != nil {
		return 0, err
	}
	sqlQuery, args, err := psql.
		Insert("posts").Columns(POST_COLUMNS[1:]...).
		Values(newPost.Title, newPost.Slug, newPost.Description, newPost.Image, newPost.Body, newPost.Published, newPost.CreatedAt, newPost.UpdatedAt,
			newPost.SocialTitle, newPost.SocialDescription, newPost.Image, newPost.Author,
			meta).Suffix("RETURNING \"id\"").

		ToSql()
	result, err := up.db.Query(sqlQuery, args...)
	if err != nil {
		if strings.Contains(err.Error(), "posts_slug_key") {
			return 0, internal_errors.ERROR_DUPLICATE
		}
		return 0, err
	}
	var lastID int
	result.Next()
	err = result.Scan(&lastID)
	if err != nil {
		return 0, err

	}
	return lastID, nil
}

func (up *PostRepository) CreateMultipleNewPost(newPosts []*model.Post) error {
	//up.db.Exec()
	var sqlBuilder = psql.Insert("posts").Columns(POST_COLUMNS[1:]...)
	for _, newPost := range newPosts {
		newPost.CreatedAt = utils.String(utils.GetCurrentISOTime())
		newPost.UpdatedAt = utils.String(utils.GetCurrentISOTime())
		meta, err := json.Marshal(newPost.Meta)
		if err != nil {
			return err
		}
		sqlBuilder = sqlBuilder.
			Values(newPost.Title, newPost.Slug, newPost.Description, newPost.Image, newPost.Body, newPost.Published, newPost.CreatedAt, newPost.UpdatedAt,
				newPost.SocialTitle, newPost.SocialDescription, newPost.Image, newPost.Author,
				meta)
	}
	sqlQuery, args, err := sqlBuilder.ToSql()
	if err != nil {
		return err
	}
	_, err = up.db.Query(sqlQuery, args...)

	return err
}

func (up *PostRepository) UpdatePost(updatePost *model.Post) error {
	builder := psql.Update("posts")
	if updatePost.Title != nil {
		builder = builder.Set("title", updatePost.Title)
	}
	if updatePost.Slug != nil {
		builder = builder.Set("slug", updatePost.Slug)
	}
	if updatePost.Description != nil {
		builder = builder.Set("description", updatePost.Description)
	}
	if updatePost.Body != nil {
		builder = builder.Set("body", updatePost.Body)
	}
	if updatePost.Published != nil {
		builder = builder.Set("published", updatePost.Published)
	}
	if updatePost.SocialImage != nil {
		builder = builder.Set("social_image", updatePost.SocialImage)
	}
	if updatePost.SocialTitle != nil {
		builder = builder.Set("social_title", updatePost.SocialTitle)
	}
	if updatePost.SocialDescription != nil {
		builder = builder.Set("social_description", updatePost.SocialDescription)
	}
	if updatePost.Author != nil {
		builder = builder.Set("author", updatePost.Author)
	}

	if updatePost.Meta != nil {
		meta, err := json.Marshal(updatePost.Meta)
		if err != nil {
			return err
		}
		builder = builder.Set("meta", meta)
	}
	sqlQuery, args, err := builder.
		Set("updated_at", utils.GetCurrentTimeStamp()).
		Where("id = ?", updatePost.ID).
		ToSql()
	_, err = up.db.Exec(sqlQuery, args...)

	if err != nil {
		if strings.Contains(err.Error(), "posts_slug_key") {
			log.Printf("%s", err)
			return internal_errors.ERROR_DUPLICATE
		}
		return err
	}
	return nil
}

func (up *PostRepository) Delete(id int) error {
	builder := psql.Update("posts")

	sqlQuery, args, err := builder.
		Set("updated_at", utils.GetCurrentTimeStamp()).
		Set("deleted", true).
		Where("id = ?", id).
		ToSql()
	_, err = up.db.Exec(sqlQuery, args...)

	if err != nil {
		if strings.Contains(err.Error(), "posts_slug_key") {
			log.Printf("%s", err)
			return internal_errors.ERROR_DUPLICATE
		}
		return err
	}
	return nil
}
