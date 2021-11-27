package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Client struct {
	db            *sqlx.DB
	connectionURL string
}

func NewClient(connectionURL string) (*Client, error) {
	c := &Client{
		connectionURL: connectionURL,
	}

	return c, nil
}

func (c *Client) Connect() error {
	db, err := sqlx.Connect("postgres", c.connectionURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	c.db = db
	return nil
}

func (c *Client) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return c.db.GetContext(ctx, dest, query, args...)
}

func (c *Client) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return c.db.SelectContext(ctx, dest, query, args...)
}

func (c *Client) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return c.db.ExecContext(ctx, query, args...)
}

func (c *Client) Close() error {
	return c.db.Close()
}
