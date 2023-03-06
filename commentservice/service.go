package commentservice

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ServiceWeaver/weaver"
	_ "github.com/lib/pq"
)

type Comment struct {
	weaver.AutoMarshal
	ID       int       `json:"id"`
	PostID   int       `json:"postId"`
	Content  string    `json:"content"`
	CreateAt time.Time `json:"createAt"`
}

func NewComment(postID int, content string) *Comment {
	return &Comment{
		PostID:   postID,
		Content:  content,
		CreateAt: time.Now().UTC(),
	}
}

type Service interface {
	GetByPost(context.Context, int) ([]Comment, error)
	CreateComment(ctx context.Context, postID int, content string) error
}

type config struct {
	URI string `toml:"db_uri"`
}

type commentService struct {
	weaver.Implements[Service]
	weaver.WithConfig[config]

	db *sql.DB
}

func (s *commentService) Init(ctx context.Context) error {
	cfg := s.Config()

	var err error
	s.db, err = sql.Open("postgres", cfg.URI)
	if err != nil {
		return fmt.Errorf("error opening database %s: %w", cfg.URI, err)
	}

	query := `
		CREATE TABLE IF NOT EXISTS comments (
			id			SERIAL PRIMARY KEY,
			post_id		INTEGER NOT NULL,
			content		TEXT NOT NULL,
			create_at 	TIMESTAMP NOT NULL
		)
	`

	_, err = s.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("error initializing database %s: %w", cfg.URI, err)
	}

	return nil
}

func (s *commentService) GetByPost(ctx context.Context, postID int) ([]Comment, error) {
	s.Logger().Info("getting comments by post id", "postID", postID)

	query := "select id, post_id, content, create_at from comments where post_id = $1 limit 5"
	rows, err := s.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var com Comment
		err := rows.Scan(&com.ID, &com.PostID, &com.Content, &com.CreateAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, com)
	}

	return comments, nil
}

func (s *commentService) CreateComment(ctx context.Context, postID int, content string) error {
	s.Logger().Info("creating comment", "postID", postID, "content", content)
	comment := NewComment(postID, content)
	query := "insert into comments ( post_id, content, create_at  ) values ($1, $2, $3)"
	_, err := s.db.ExecContext(ctx, query, comment.PostID, comment.Content, comment.CreateAt)
	if err != nil {
		return err
	}
	return nil
}
