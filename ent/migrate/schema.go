// Code generated by entc, DO NOT EDIT.

package migrate

import (
	"github.com/facebook/ent/dialect/entsql"
	"github.com/facebook/ent/dialect/sql/schema"
	"github.com/facebook/ent/schema/field"
)

var (
	// MediaColumns holds the columns for the "media" table.
	MediaColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID, Unique: true},
		{Name: "title", Type: field.TypeString, Size: 255},
		{Name: "status", Type: field.TypeEnum, Enums: []string{"AwaitingUpload", "Processing", "Ready", "Errored"}},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
	}
	// MediaTable holds the schema information for the "media" table.
	MediaTable = &schema.Table{
		Name:        "media",
		Columns:     MediaColumns,
		PrimaryKey:  []*schema.Column{MediaColumns[0]},
		ForeignKeys: []*schema.ForeignKey{},
		Annotation:  &entsql.Annotation{Table: "media"},
	}
	// MediaFileColumns holds the columns for the "media_file" table.
	MediaFileColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID, Unique: true},
		{Name: "video_bitrate", Type: field.TypeInt64},
		{Name: "scaled_width", Type: field.TypeInt16},
		{Name: "encoder_preset", Type: field.TypeEnum, Enums: []string{"source", "ultrafast", "veryfast", "fast", "medium", "slow", "veryslow"}},
		{Name: "framerate", Type: field.TypeInt8},
		{Name: "duration_seconds", Type: field.TypeFloat64},
		{Name: "media_type", Type: field.TypeEnum, Enums: []string{"audio", "video"}},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "media", Type: field.TypeUUID, Nullable: true},
	}
	// MediaFileTable holds the schema information for the "media_file" table.
	MediaFileTable = &schema.Table{
		Name:       "media_file",
		Columns:    MediaFileColumns,
		PrimaryKey: []*schema.Column{MediaFileColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:  "media_file_media_files",
				Columns: []*schema.Column{MediaFileColumns[9]},

				RefColumns: []*schema.Column{MediaColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
		Annotation: &entsql.Annotation{Table: "media_file"},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		MediaTable,
		MediaFileTable,
	}
)

func init() {
	MediaFileTable.ForeignKeys[0].RefTable = MediaTable
}
