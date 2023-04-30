package content

import (
	"fmt"
	"github.com/monstrum/ponzu-cms/management/editor"
	"github.com/monstrum/ponzu-cms/system/item"
)

type Song struct {
	item.Item

	Title      string `json:"title"`
	Artist     string `json:"artist"`
	Rating     int    `json:"rating"`
	Opinion    string `json:"opinion"`
	SpotifyUrl string `json:"spotify_url"`
}

// MarshalEditor writes a buffer of html to edit a Song within the CMS
// and implements editor.Editable
func (s *Song) MarshalEditor() ([]byte, error) {
	view, err := editor.Form(s,
		// Take note that the first argument to these Input-like functions
		// is the string version of each Song field, and must follow
		// this pattern for auto-decoding and auto-encoding reasons:
		editor.Field{
			View: editor.Input("Title", s, map[string]string{
				"label":       "Title",
				"type":        "text",
				"placeholder": "Enter the Title here",
			}),
		},
		editor.Field{
			View: editor.Input("Artist", s, map[string]string{
				"label":       "Artist",
				"type":        "text",
				"placeholder": "Enter the Artist here",
			}),
		},
		editor.Field{
			View: editor.Input("Rating", s, map[string]string{
				"label":       "Rating",
				"type":        "text",
				"placeholder": "Enter the Rating here",
			}),
		},
		editor.Field{
			View: editor.Richtext("Opinion", s, map[string]string{
				"label":       "Opinion",
				"placeholder": "Enter the Opinion here",
			}),
		},
		editor.Field{
			View: editor.Input("SpotifyUrl", s, map[string]string{
				"label":       "SpotifyUrl",
				"type":        "text",
				"placeholder": "Enter the SpotifyUrl here",
			}),
		},
	)

	if err != nil {
		return nil, fmt.Errorf("Failed to render Song editor view: %s", err.Error())
	}

	return view, nil
}

func init() {
	item.Types["Song"] = func() interface{} { return new(Song) }
}

// String defines how a Song is printed. Update it using more descriptive
// fields from the Song struct type
func (s *Song) String() string {
	return fmt.Sprintf("Song: %s", s.UUID)
}
