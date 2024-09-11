package models

import (
	"context"
	"time"

	"github.com/charmbracelet/log"
	"github.com/nathanhollows/scanner-server/db"
)

type Scan struct {
	LocationID string    `bun:",notnull" json:"location_id"`
	TagID      string    `bun:",notnull" json:"tag"`
	Timestamp  time.Time `bun:",notnull" json:"timestamp"`
}

type Scans []Scan

// Find scans by Tag
func FindScansByTag(ctx context.Context, list_id string) Scans {
	var tag Tag
	err := db.DB.NewSelect().
		Model(&tag).
		Where("list_id = ?", list_id).
		Scan(ctx)
	if err != nil {
		log.Error("error finding tag", "err", err, "tag", list_id)
		return Scans{}
	}
	var scans Scans
	err = db.DB.NewSelect().
		Model(&scans).
		ColumnExpr("max(timestamp) as timestamp, location_id, tag_id").
		Where("tag_id = ?", tag.TagID).
		Group("location_id").
		Order("timestamp DESC").
		Scan(ctx)
	if err != nil {
		log.Error("error finding scans", "err", err, "tag", list_id)
		return Scans{}
	}
	return scans
}

// Create a new scan
func (scan *Scan) Save(ctx context.Context) error {
	_, err := db.DB.NewInsert().Model(scan).Exec(ctx)
	return err
}
