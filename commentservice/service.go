package commentservice

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ServiceWeaver/weaver"
	_ "github.com/lib/pq"
)

type Comments struct {
	weaver.AutoMarshal
	ID       int
	PostID   int
	Content  string
	CreateAt time.Time
}

type Service interface {
	GetByPost(context.Context, int) ([]Comments, error)
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

func (s *commentService) GetByPost(ctx context.Context, postID int) ([]Comments, error) {
	s.Logger().Info("getting comments by post id", "postID", postID)

	rows, err := s.db.QueryContext(ctx, "select id, post_id, content, create_at from comments where post_id = $1", postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comments
	for rows.Next() {
		var com Comments
		err := rows.Scan(&com.ID, &com.PostID, &com.Content, &com.CreateAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, com)
	}

	return comments, nil
}
