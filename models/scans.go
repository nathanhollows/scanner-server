package models

import (
	"context"
	"time"

	"github.com/nathanhollows/scanner-server/db"
)

type Scan struct {
	LocationID string    `bun:",notnull" json:"location_id"`
	TagID      string    `bun:",notnull" json:"tag"`
	Timestamp  time.Time `bun:",notnull" json:"timestamp"`
}

type Scans []Scan

// Find scans by Tag
func FindScansByTag(ctx context.Context, tag string) Scans {
	var scans Scans
	db.DB.NewSelect().
		Model(&scans).
		ColumnExpr("max(timestamp) as timestamp, location_id, tag_id").
		Group("location_id").
		Order("timestamp DESC").
		Scan(ctx)
	return scans
}

// Create a new scan
func (scan *Scan) Save(ctx context.Context) error {
	_, err := db.DB.NewInsert().Model(scan).Exec(ctx)
	return err
}
