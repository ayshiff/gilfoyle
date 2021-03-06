// Code generated by entc, DO NOT EDIT.

package mediafile

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the mediafile type in the database.
	Label = "media_file"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldVideoBitrate holds the string denoting the video_bitrate field in the database.
	FieldVideoBitrate = "video_bitrate"
	// FieldScaledWidth holds the string denoting the scaled_width field in the database.
	FieldScaledWidth = "scaled_width"
	// FieldEncoderPreset holds the string denoting the encoder_preset field in the database.
	FieldEncoderPreset = "encoder_preset"
	// FieldFramerate holds the string denoting the framerate field in the database.
	FieldFramerate = "framerate"
	// FieldDurationSeconds holds the string denoting the duration_seconds field in the database.
	FieldDurationSeconds = "duration_seconds"
	// FieldMediaType holds the string denoting the media_type field in the database.
	FieldMediaType = "media_type"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"

	// EdgeMedia holds the string denoting the media edge name in mutations.
	EdgeMedia = "media"

	// Table holds the table name of the mediafile in the database.
	Table = "media_file"
	// MediaTable is the table the holds the media relation/edge.
	MediaTable = "media_file"
	// MediaInverseTable is the table name for the Media entity.
	// It exists in this package in order to avoid circular dependency with the "media" package.
	MediaInverseTable = "media"
	// MediaColumn is the table column denoting the media relation/edge.
	MediaColumn = "media"
)

// Columns holds all SQL columns for mediafile fields.
var Columns = []string{
	FieldID,
	FieldVideoBitrate,
	FieldScaledWidth,
	FieldEncoderPreset,
	FieldFramerate,
	FieldDurationSeconds,
	FieldMediaType,
	FieldCreatedAt,
	FieldUpdatedAt,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the MediaFile type.
var ForeignKeys = []string{
	"media",
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	for i := range ForeignKeys {
		if column == ForeignKeys[i] {
			return true
		}
	}
	return false
}

var (
	// VideoBitrateValidator is a validator for the "video_bitrate" field. It is called by the builders before save.
	VideoBitrateValidator func(int64) error
	// ScaledWidthValidator is a validator for the "scaled_width" field. It is called by the builders before save.
	ScaledWidthValidator func(int16) error
	// FramerateValidator is a validator for the "framerate" field. It is called by the builders before save.
	FramerateValidator func(int8) error
	// DurationSecondsValidator is a validator for the "duration_seconds" field. It is called by the builders before save.
	DurationSecondsValidator func(float64) error
	// DefaultCreatedAt holds the default value on creation for the created_at field.
	DefaultCreatedAt func() time.Time
	// DefaultUpdatedAt holds the default value on creation for the updated_at field.
	DefaultUpdatedAt func() time.Time
	// UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	UpdateDefaultUpdatedAt func() time.Time
	// DefaultID holds the default value on creation for the id field.
	DefaultID func() uuid.UUID
)

// EncoderPreset defines the type for the encoder_preset enum field.
type EncoderPreset string

// EncoderPreset values.
const (
	EncoderPresetSource    EncoderPreset = "source"
	EncoderPresetUltrafast EncoderPreset = "ultrafast"
	EncoderPresetVeryfast  EncoderPreset = "veryfast"
	EncoderPresetFast      EncoderPreset = "fast"
	EncoderPresetMedium    EncoderPreset = "medium"
	EncoderPresetSlow      EncoderPreset = "slow"
	EncoderPresetVeryslow  EncoderPreset = "veryslow"
)

func (ep EncoderPreset) String() string {
	return string(ep)
}

// EncoderPresetValidator is a validator for the "encoder_preset" field enum values. It is called by the builders before save.
func EncoderPresetValidator(ep EncoderPreset) error {
	switch ep {
	case EncoderPresetSource, EncoderPresetUltrafast, EncoderPresetVeryfast, EncoderPresetFast, EncoderPresetMedium, EncoderPresetSlow, EncoderPresetVeryslow:
		return nil
	default:
		return fmt.Errorf("mediafile: invalid enum value for encoder_preset field: %q", ep)
	}
}

// MediaType defines the type for the media_type enum field.
type MediaType string

// MediaType values.
const (
	MediaTypeAudio MediaType = "audio"
	MediaTypeVideo MediaType = "video"
)

func (mt MediaType) String() string {
	return string(mt)
}

// MediaTypeValidator is a validator for the "media_type" field enum values. It is called by the builders before save.
func MediaTypeValidator(mt MediaType) error {
	switch mt {
	case MediaTypeAudio, MediaTypeVideo:
		return nil
	default:
		return fmt.Errorf("mediafile: invalid enum value for media_type field: %q", mt)
	}
}
