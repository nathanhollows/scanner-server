package models

import (
	"context"

	"github.com/nathanhollows/scanner-server/db"
)

type Tag struct {
	TagID  string `bun:",pk"`
	ListID string `bun:",pk"`
}

type Tags []Tag

// Save saves the tag to the database.
func (t *Tag) Save(ctx context.Context) error {
	_, err := db.DB.NewInsert().Model(t).Exec(ctx)
	return err
}
