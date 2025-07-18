package feedreader

import (
	"testing"

	"github.com/mmcdole/gofeed"
	"gotest.tools/assert"
)

func Test_getShowLines(t *testing.T) {
	tests := []struct {
		name      string
		feedItem  *FeedItem
		showType  ShowType
		expected1 string
		expected2 string
	}{
		{
			name:      "with nil FeedItem",
			feedItem:  nil,
			showType:  SHOW_TITLE,
			expected1: "",
			expected2: "",
		},
		{
			name: "with plain title",
			feedItem: &FeedItem{
				item: &gofeed.Item{Title: "Cats and Dogs", Link: "https://cats.com"},
			},
			showType:  SHOW_TITLE,
			expected1: "[white]Cats and Dogs",
			expected2: "https://cats.com",
		},
		{
			name: "with link",
			feedItem: &FeedItem{
				item: &gofeed.Item{Title: "Cats and Dogs", Link: "https://cats.com/dog.xml"},
			},
			showType:  SHOW_LINK,
			expected1: "https://cats.com/dog.xml",
			expected2: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			widget := &Widget{
				settings: &Settings{
					colors: colors{
						source:      "green",
						publishDate: "orange",
					},
					showSource: true,
				},
				showType: tt.showType,
			}

			l1, l2 := widget.getShowLines(tt.feedItem, "white")

			assert.Equal(t, tt.expected1, l1)
			assert.Equal(t, tt.expected2, l2)
		})
	}
}
