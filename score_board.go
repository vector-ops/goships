package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rthornton128/goncurses"
	"github.com/vector-ops/goships/types"
)

type ScoreBoard struct {
	win *goncurses.Window

	title        string
	titleColor   int16
	CurrentScore Score
	OverallScore Score
}

type Score struct {
	PlayerScore int
	EnemyScore  int
}

func NewScoreBoard(win *goncurses.Window) *ScoreBoard {
	return &ScoreBoard{
		win:        win,
		title:      "SCORE",
		titleColor: types.BLUE_BLACK,
	}
}

func (s *ScoreBoard) Render(ctx context.Context) error {
	s.win.Erase()
	err := s.win.Box(goncurses.ACS_VLINE, goncurses.ACS_HLINE)
	if err != nil {
		return err
	}
	// _, mx := s.win.MaxYX()

	// s.win.ColorOn(s.titleColor)
	// s.win.MovePrint(1, (mx/2)-len(s.title)/2, s.title)
	// s.win.ColorOff(s.titleColor)
	s.draw()
	s.win.NoutRefresh()

	return nil
}

func (s *ScoreBoard) Close() error {
	return s.win.Delete()
}

func (s *ScoreBoard) draw() error {
	startX := 2
	startY := 1
	titleOffset := 2

	// draw score board
	scoreStat := [][]string{
		{"Player", "Enemy"},
		{strconv.Itoa(s.CurrentScore.PlayerScore), strconv.Itoa(s.CurrentScore.EnemyScore)},
	}
	s.drawStatBoard(startX, startY, "SCORE", scoreStat)
	startY += (len(scoreStat[0]) + titleOffset) * 2

	// draw alive ships (later reuse for showing available ships at the beginning of the game)
	shipStat := [][]string{
		{"Player", "Enemy"},
		{"3", "4"},
	}
	s.drawStatBoard(startX, startY, "SHIPS", shipStat)

	return nil
}

func (s *ScoreBoard) drawStatBoard(startX, startY int, title string, stats [][]string) error {
	my, mx := s.win.MaxYX()

	if len(stats) == 0 {
		return nil
	}

	if len(stats[0]) == 0 {
		return nil
	}

	if startX < 0 || startY < 0 {
		return fmt.Errorf("stat board %s: invalid position %d, %d", title, startX, startY)
	}

	rows := len(stats)
	cols := len(stats[0])

	if startX+cols > mx || startY+rows > my {
		return fmt.Errorf("stat board %s: position out of bounds %d, %d", title, startX, startY)
	}

	maxColWidth := ((mx - startX) / cols) - 1
	maxRowHeight := 2

	// draw title row upper border
	s.win.ColorOn(types.COLOR_WALL)
	s.win.MoveAddChar(startY, startX, goncurses.Char(types.WALLS_ASCII[types.CELL_WALL_CORNER]))
	for i := 1; i <= (maxColWidth * cols); i++ {
		x := startX + i
		y := startY

		s.win.MoveAddChar(y, x, goncurses.Char(types.WALLS_ASCII[types.CELL_WALL_HORIZONTAL]))
		if i%(maxColWidth*cols) == 0 {
			s.win.MoveAddChar(y, x, goncurses.Char(types.WALLS_ASCII[types.CELL_WALL_CORNER]))
		}
	}
	s.win.ColorOff(types.COLOR_WALL)
	startY++

	// print title and draw vertical lines enclosing the title
	for i := 0; i <= (maxColWidth * cols); i++ {
		x := startX + i
		y := startY

		s.win.ColorOn(types.GREEN_BLACK)
		titleX := ((maxColWidth*cols)-len(title)/2)/2 + 1
		if i == titleX {
			s.win.MovePrint(y, titleX, title)
		}
		s.win.ColorOff(types.GREEN_BLACK)

		s.win.ColorOn(types.COLOR_WALL)
		if i%(maxColWidth*cols) == 0 {
			s.win.MoveAddChar(y, x, goncurses.Char(types.WALLS_ASCII[types.CELL_WALL_VERTICAL]))
		}
		s.win.ColorOff(types.COLOR_WALL)
	}
	startY++

	for row := 0; row <= rows; row++ {
		for col := 0; col <= cols; col++ {

			y := startY + row*maxRowHeight
			x := startX + col*maxColWidth

			s.win.ColorOn(types.COLOR_WALL)
			s.win.MoveAddChar(y, x, goncurses.Char(types.WALLS_ASCII[types.CELL_WALL_CORNER]))

			if col < cols {
				for i := 1; i < maxColWidth; i++ {
					s.win.MoveAddChar(y, x+i, goncurses.Char(types.WALLS_ASCII[types.CELL_WALL_HORIZONTAL]))
				}
			}

			if row < rows {
				s.win.MoveAddChar(y+1, x, goncurses.Char(types.WALLS_ASCII[types.CELL_WALL_VERTICAL]))
			}
			s.win.ColorOff(types.COLOR_WALL)

			if row == 0 { // blue black for the header
				s.win.ColorOn(types.BLUE_BLACK)
			} else { // white black for the rest
				s.win.ColorOn(types.WHITE_BLACK)
			}

			if row < rows && col < cols {
				contentX := x + maxColWidth/2 - len(stats[row][col])/2
				s.win.MovePrint(y+1, contentX, stats[row][col])
			}

			if row == 0 {
				s.win.ColorOff(types.BLUE_BLACK)
			} else {
				s.win.ColorOff(types.WHITE_BLACK)
			}

		}
	}

	return nil
}
