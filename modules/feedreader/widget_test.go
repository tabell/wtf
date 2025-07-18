package feedreader

import (
	"fmt"
	"strings"
	"testing"

	"github.com/mmcdole/gofeed"
	"github.com/rivo/tview"
	"github.com/wtfutil/wtf/cfg"
	"github.com/wtfutil/wtf/utils"
	"gotest.tools/assert"
)

func Test_getShowText(t *testing.T) {
	tests := []struct {
		name     string
		feedItem *FeedItem
		showType ShowType
		expected string
	}{
		{
			name:     "with nil FeedItem",
			feedItem: nil,
			showType: SHOW_TITLE,
			expected: "",
		},
		{
			name: "with plain title",
			feedItem: &FeedItem{
				item: &gofeed.Item{Title: "Cats and Dogs"},
			},
			showType: SHOW_TITLE,
			expected: "[white]Cats and Dogs",
		},
		{
			name: "with escaped title",
			feedItem: &FeedItem{
				item: &gofeed.Item{Title: "&lt;Cats and Dogs&gt;"},
			},
			showType: SHOW_TITLE,
			expected: "[white]<Cats and Dogs>",
		},
		{
			name: "with unescaped title",
			feedItem: &FeedItem{
				item: &gofeed.Item{Title: "<Cats and Dogs>"},
			},
			showType: SHOW_TITLE,
			expected: "[white]<Cats and Dogs>",
		},
		{
			name: "with source-title",
			feedItem: &FeedItem{
				sourceTitle: "WTF",
				item:        &gofeed.Item{Title: "<Cats and Dogs>"},
			},
			showType: SHOW_TITLE,
			expected: "[green]WTF [white]<Cats and Dogs>",
		},
		{
			name: "with link",
			feedItem: &FeedItem{
				item: &gofeed.Item{Title: "Cats and Dogs", Link: "https://cats.com/dog.xml"},
			},
			showType: SHOW_LINK,
			expected: "https://cats.com/dog.xml",
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

			actual := widget.getShowText(tt.feedItem, "white")

			assert.Equal(t, tt.expected, actual)
		})
	}
}

func Test_widget_content_block(t *testing.T) {
	app := tview.NewApplication()
	w := NewWidget(app, make(chan bool), tview.NewPages(), &Settings{Common: &cfg.Common{}, maxHeight: 3})
	w.showType = SHOW_CONTENT
	w.stories = []*FeedItem{
		{
			item: &gofeed.Item{
				Title:   "Cats",
				Content: "<pre>one\ntwo\nthree\nfour</pre>",
			},
		},
	}

	title, content, wrap := w.content()

	rowColor := w.RowColor(0)
	display := w.getShowText(w.stories[0], rowColor)
	lines := strings.Split(display, "\n")
	if w.settings.maxHeight == 0 || w.settings.maxHeight >= 2 {
		lines = append([]string{""}, lines...)
	}
	if w.settings.maxHeight > 0 && len(lines) > w.settings.maxHeight {
		lines = lines[:w.settings.maxHeight]
	}
	if w.settings.maxHeight == 0 || w.settings.maxHeight >= 2 {
		lines[0] = fmt.Sprintf("[%s]%2d.[white]", rowColor, 1)
		for i := 1; i < len(lines); i++ {
			lines[i] = fmt.Sprintf("[%s]%s[white]", rowColor, lines[i])
		}
	} else {
		lines[0] = fmt.Sprintf("[%s]%2d. %s[white]", rowColor, 1, lines[0])
		for i := 1; i < len(lines); i++ {
			lines[i] = fmt.Sprintf("[%s]%s[white]", rowColor, lines[i])
		}
	}
	expected := utils.HighlightableBlockHelper(w.View, lines, 0)

	assert.Equal(t, w.CommonSettings().Title, title)
	assert.Equal(t, expected, content)
	assert.Equal(t, true, wrap)
}

func Test_widget_content_min_height(t *testing.T) {
	app := tview.NewApplication()
	w := NewWidget(app, make(chan bool), tview.NewPages(), &Settings{Common: &cfg.Common{}, minHeight: 2})
	w.showType = SHOW_TITLE
	w.stories = []*FeedItem{
		{
			item: &gofeed.Item{
				Title: "Cats",
			},
		},
	}

	_, content, _ := w.content()

	expectedLines := 2
	actualLines := strings.Count(content, "\n")
	assert.Equal(t, expectedLines, actualLines)
}

func Test_widget_title_second_line(t *testing.T) {
	app := tview.NewApplication()
	w := NewWidget(app, make(chan bool), tview.NewPages(), &Settings{Common: &cfg.Common{}, maxHeight: 3})
	w.showType = SHOW_CONTENT
	w.stories = []*FeedItem{
		{
			item: &gofeed.Item{
				Title:   "Cats",
				Content: "<pre>meow</pre>",
			},
		},
	}

	_, content, _ := w.content()

	// strip highlight regions and padding for easier comparison
	cleaned := utils.StripColorTags(strings.ReplaceAll(content, "[\"0\"][\"\"]", ""))
	lines := strings.Split(cleaned, "\n")

	if len(lines) < 2 {
		t.Fatalf("expected at least 2 lines, got %d", len(lines))
	}

	expectedTitle := "Cats           "
	assert.Equal(t, expectedTitle, lines[1])
}
