package game

import (
	"context"
	"fmt"

	"github.com/rthornton128/goncurses"
	"github.com/vector-ops/goships/types"
)

type ScoreBoard struct {
	win   *goncurses.Window
	debug bool

	title        string
	titleColor   int16
	CurrentScore Score
	OverallScore Score

	stats []StatBoard
}

type Score struct {
	PlayerScore int
	EnemyScore  int
}

type StatBoard struct {
	Title      string
	StatHeader []string
	StatValues []string
}

func NewScoreBoard(win *goncurses.Window, stats []StatBoard, debug bool) *ScoreBoard {

	return &ScoreBoard{
		win:        win,
		debug:      debug,
		title:      "SCORE",
		titleColor: types.BlueBlack,
		stats:      stats,
	}
}

func (s *ScoreBoard) Render(ctx context.Context) error {
	s.win.Erase()
	err := s.win.Box(goncurses.ACS_VLINE, goncurses.ACS_HLINE)
	if err != nil {
		return err
	}

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

	for _, stat := range s.stats {
		stats := [][]string{
			stat.StatHeader,
			stat.StatValues,
		}
		s.drawStatBoard(startX, startY, stat.Title, stats)
		startY += (len(stats[0]) + titleOffset) * 2
	}

	return nil
}

func (s *ScoreBoard) SetScorePlayerScore(score int) {
	s.CurrentScore.PlayerScore = score
}

func (s *ScoreBoard) SetScoreEnemyScore(score int) {
	s.CurrentScore.EnemyScore = score
}

func (s *ScoreBoard) SetStat(title string, statValues []string) {
	for i, stat := range s.stats {
		if stat.Title == title {
			s.stats[i].StatValues = statValues
		}
	}
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
	s.win.ColorOn(types.ColorWall)
	s.win.MoveAddChar(startY, startX, goncurses.Char(types.WallsASCII[types.CellWallCorner]))
	for i := 1; i <= (maxColWidth * cols); i++ {
		x := startX + i
		y := startY

		s.win.MoveAddChar(y, x, goncurses.Char(types.WallsASCII[types.CellWallHorizontal]))
		if i%(maxColWidth*cols) == 0 {
			s.win.MoveAddChar(y, x, goncurses.Char(types.WallsASCII[types.CellWallCorner]))
		}
	}
	s.win.ColorOff(types.ColorWall)
	startY++

	// print title and draw vertical lines enclosing the title
	for i := 0; i <= (maxColWidth * cols); i++ {
		x := startX + i
		y := startY

		s.win.ColorOn(types.GreenBlack)
		titleX := ((maxColWidth*cols)-len(title)/2)/2 + 1
		if i == titleX {
			s.win.MovePrint(y, titleX, title)
		}
		s.win.ColorOff(types.GreenBlack)

		s.win.ColorOn(types.ColorWall)
		if i%(maxColWidth*cols) == 0 {
			s.win.MoveAddChar(y, x, goncurses.Char(types.WallsASCII[types.CellWallVertical]))
		}
		s.win.ColorOff(types.ColorWall)
	}
	startY++

	for row := 0; row <= rows; row++ {
		for col := 0; col <= cols; col++ {

			y := startY + row*maxRowHeight
			x := startX + col*maxColWidth

			s.win.ColorOn(types.ColorWall)
			s.win.MoveAddChar(y, x, goncurses.Char(types.WallsASCII[types.CellWallCorner]))

			if col < cols {
				for i := 1; i < maxColWidth; i++ {
					s.win.MoveAddChar(y, x+i, goncurses.Char(types.WallsASCII[types.CellWallHorizontal]))
				}
			}

			if row < rows {
				s.win.MoveAddChar(y+1, x, goncurses.Char(types.WallsASCII[types.CellWallVertical]))
			}
			s.win.ColorOff(types.ColorWall)

			if row == 0 { // blue black for the header
				s.win.ColorOn(types.BlueBlack)
			} else { // white black for the rest
				s.win.ColorOn(types.WhiteBlack)
			}

			if row < rows && col < cols {
				contentX := x + maxColWidth/2 - len(stats[row][col])/2
				s.win.MovePrint(y+1, contentX, stats[row][col])
			}

			if row == 0 {
				s.win.ColorOff(types.BlueBlack)
			} else {
				s.win.ColorOff(types.WhiteBlack)
			}

		}
	}

	return nil
}
