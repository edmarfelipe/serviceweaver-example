package postservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/ServiceWeaver/weaver"
	_ "github.com/lib/pq"
)

func createSlug(s string) string {
	var re = regexp.MustCompile("[^a-z0-9]+")
	return strings.Trim(re.ReplaceAllString(strings.ToLower(s), "-"), "-")
}

type Post struct {
	weaver.AutoMarshal
	ID       int       `json:"id"`
	Title    string    `json:"title"`
	Slug     string    `json:"slug"`
	Content  string    `json:"content"`
	CreateAt time.Time `json:"createAt"`
}

func NewPost(title string, content string) *Post {
	return &Post{
		Title:    title,
		Slug:     createSlug(title),
		Content:  content,
		CreateAt: time.Now().UTC(),
	}
}

type Service interface {
	CreatePost(context.Context, string, string) error
	GetPost(context.Context, string) (*Post, error)
	GetLatestPosts(context.Context, int) ([]Post, error)
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
			slug		VARCHAR(100) NOT NULL UNIQUE,
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

func (s *postService) GetLatestPosts(ctx context.Context, offset int) ([]Post, error) {
	s.Logger().Info("getting latest posts", "offset", offset)
	query := "select id, title, slug, content, create_at from posts order by create_at limit 10 offset $1"
	rows, err := s.db.QueryContext(ctx, query, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Title, &post.Slug, &post.Content, &post.CreateAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (s *postService) CreatePost(ctx context.Context, title string, content string) error {
	s.Logger().Info("creating post", "title", title, "content", content)
	post := NewPost(title, content)
	query := "insert into posts (title, slug, content, create_at) values ($1, $2, $3, $4)"
	_, err := s.db.ExecContext(ctx, query, post.Title, post.Slug, post.Content, post.CreateAt)
	if err != nil {
		return err
	}
	return nil
}
