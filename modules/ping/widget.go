package ping

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	probing "github.com/prometheus-community/pro-bing"
	"github.com/rivo/tview"
	"github.com/wtfutil/wtf/view"
)

// Widget is the container for your module's data
type Widget struct {
	view.TextWidget
	hosts []Host

	settings *Settings
}

// NewWidget creates and returns an instance of Widget
func NewWidget(tviewApp *tview.Application, redrawChan chan bool, settings *Settings) *Widget {
	widget := Widget{
		TextWidget: view.NewTextWidget(tviewApp, redrawChan, nil, settings.common),

		settings: settings,
	}
	widget.hosts = widget.settings.hosts

	return &widget
}

/* -------------------- Exported Functions -------------------- */

func (widget *Widget) doPings() {
	var wg sync.WaitGroup
	for i := range widget.hosts {
		idx := i
		host := widget.hosts[idx]
		widget.hosts[idx].Up = false // reset to false each time
		wg.Add(1)
		go func() {
			defer wg.Done()
			pinger, err := probing.NewPinger(host.Hostname)
			if err == nil {
				pinger.Count = 1
				pinger.Timeout = 10 * time.Second
				err = pinger.Run() // Blocks until finished.
				if err == nil {
					stats := pinger.Statistics() // get send/receive/duplicate/rtt stats
					if stats.PacketsRecv > 0 {
						widget.hosts[idx].Up = true
					} else {
						widget.hosts[idx].Up = false
					}
				} else {
					log.Fatalf("error sending ping: %v", err)
				}
			}

		}()
	}
	wg.Wait()
}
func (widget *Widget) Refresh() {

	widget.doPings()
	widget.display()
}

/* -------------------- Unexported Functions -------------------- */

func (widget *Widget) content() string {
	nameWidth := 12
	for _, t := range widget.hosts {
		if len(t.Label) > nameWidth {
			nameWidth = len(t.Label) + 2
		}
	}

	s := []string{}
	for _, t := range widget.hosts {
		var status string
		if t.Up {
			status = "[green]Up"
		} else {
			status = "[red]DOWN"
		}
		statusLine := fmt.Sprintf("[white]%-*s: %s", nameWidth, t.Label, status)
		s = append(s, statusLine)
	}

	return strings.Join(s, "\n")
}

func (widget *Widget) display() {
	widget.Redraw(func() (string, string, bool) {
		return widget.CommonSettings().Title, widget.content(), false
	})
}
