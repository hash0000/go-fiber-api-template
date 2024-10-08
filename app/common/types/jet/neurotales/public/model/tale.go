//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import (
	"time"
)

type Tale struct {
	ID                   int64 `sql:"primary_key"`
	UserID               int64
	Name                 string
	FileName             string
	IsPayed              bool
	TaleGenerationID     string
	CreatedAt            time.Time
	ChildData            string
	BackgroundCharacters string
	Preferences          string
	Moral                string
	OpenAiAnswer         string
	FabulaImgToTextJSON  *string
}
