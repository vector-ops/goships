package main

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/rthornton128/goncurses"
	"github.com/vector-ops/goships/logger"
	"github.com/vector-ops/goships/types"
)

type LogWindow struct {
	win *goncurses.Window

	title      string
	titleColor int16

	maxLogs int

	gamelogs []logger.Log

	logCh   chan logger.Log
	closeCh chan struct{}
}

func NewLogWindow(win *goncurses.Window, logCh chan logger.Log) *LogWindow {
	my, _ := win.MaxYX()

	maxLogs := my - 8

	lw := &LogWindow{
		win:        win,
		title:      "DEBUG LOG",
		titleColor: types.BlueBlack,
		maxLogs:    maxLogs,
		gamelogs:   make([]logger.Log, maxLogs),
		logCh:      logCh,
		closeCh:    make(chan struct{}),
	}

	go lw.monitorLogs()

	return lw
}

func (l *LogWindow) Render(ctx context.Context) error {
	l.win.Erase()
	err := l.win.Box(goncurses.ACS_VLINE, goncurses.ACS_HLINE)
	if err != nil {
		return err
	}

	_, mx := l.win.MaxYX()

	l.win.ColorOn(l.titleColor)
	l.win.MovePrint(1, (mx/2)-len(l.title)/2, l.title)
	l.win.ColorOff(l.titleColor)

	l.draw()

	l.win.NoutRefresh()

	return nil
}

func (l *LogWindow) Close() error {
	close(l.closeCh)
	return l.win.Delete()
}

func (l *LogWindow) draw() {

	offsetY := 4

	my, mx := l.win.MaxYX()
	seq := 1
	for i, lg := range l.gamelogs {

		if lg.Level == "" || lg.Msg == "" || lg.Timestamp.IsZero() {
			offsetY -= 1
			continue
		}

		parts := splitStrWidth(fmt.Sprintf("[%d] %s : %s", seq, lg.Level.Upper(), lg.Msg), my, mx, 4, 6)
		l.win.ColorOn(lg.Color)
		for j, s := range parts {
			offsetY += j
			l.win.MovePrint(i+offsetY, 4, s)
		}
		l.win.ColorOff(lg.Color)
		seq++
	}
}

func splitStrWidth(str string, my, mx, offY, offX int) []string {

	parts := make([]string, 0)

	lim := mx - (2 * offX)
	if len(str) > lim {
		words := strings.Split(str, " ")

		l := 0
		fits := make([]string, 0)
		for i, w := range words {

			l += len(w)
			if l+len(w) < lim {
				fits = append(fits, w)
			} else {
				parts = append(parts, strings.Join(fits, " "))
				fits = fits[:0]
				fits = append(fits, w)
				l = len(w)
			}

			if i == len(words)-1 {
				parts = append(parts, strings.Join(fits, " "))
			}
		}

	} else {
		parts = append(parts, str)
	}

	return parts
}

func (l *LogWindow) monitorLogs() {

	for {
		select {
		case <-l.closeCh:
			return
		case log := <-l.logCh:
			if log.Level.String() == "" || log.Msg == "" || log.Timestamp.IsZero() {
				continue
			}

			l.gamelogs = append(l.gamelogs, log)
			l.rotateLogs()
		}
	}
}

func (l *LogWindow) rotateLogs() {
	if len(l.gamelogs) < l.maxLogs {
		return
	}

	i := len(l.gamelogs) - l.maxLogs

	l.gamelogs = l.gamelogs[i:]

	slices.SortFunc(l.gamelogs, func(a logger.Log, b logger.Log) int {
		t1 := a.Timestamp.UnixNano()
		t2 := b.Timestamp.UnixNano()

		if t1 < t2 {
			return -1
		} else if t1 > t2 {
			return 1
		} else {
			return 0
		}
	})

}
