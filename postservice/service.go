package postservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ServiceWeaver/weaver"
	_ "github.com/lib/pq"
)

type Post struct {
	weaver.AutoMarshal
	ID       int
	Title    string
	Slug     string
	Content  string
	CreateAt time.Time
}

type Service interface {
	GetPost(context.Context, string) (*Post, error)
}

type config struct {
	URI string `toml:"db_uri"`
}

type postService struct {
	weaver.Implements[Service]
	weaver.WithConfig[config]

	db *sql.DB
}

func (s *postService) Init(ctx context.Context) error {
	cfg := s.Config()

	var err error
	s.db, err = sql.Open("postgres", cfg.URI)
	if err != nil {
		return fmt.Errorf("error opening database %s: %w", cfg.URI, err)
	}

	query := `
		CREATE TABLE IF NOT EXISTS posts (
			id			SERIAL PRIMARY KEY,
			title		VARCHAR(100) NOT NULL,
			slug		VARCHAR(100) NOT NULL,
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

func (s *postService) GetPost(ctx context.Context, slug string) (*Post, error) {
	s.Logger().Info("getting post by slug", "slug", slug)
	post := Post{}
	query := "select id, title, slug, content, create_at from posts where slug = $1"
	err := s.db.QueryRowContext(ctx, query, slug).Scan(&post.ID, &post.Title, &post.Slug, &post.Content, &post.CreateAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &post, nil
}
