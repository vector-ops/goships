package main

import (
	"context"
	"fmt"

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
	_, mx := s.win.MaxYX()

	s.win.ColorOn(s.titleColor)
	s.win.MovePrint(1, (mx/2)-len(s.title)/2, s.title)
	s.win.ColorOff(s.titleColor)
	s.draw()
	s.win.NoutRefresh()

	return nil
}

func (s *ScoreBoard) Close() error {
	return s.win.Delete()
}

func (s *ScoreBoard) draw() error {
	// s.drawScoreBoard()
	s.drawStatBoard(2, 3, "Score Board", [][]string{
		{"Player", "Enemy"},
		{"3", "4"},
	})
	s.drawAliveShips()
	s.win.ColorOn(types.WHITE_BLACK)
	s.win.MovePrint(7, 2, fmt.Sprintf("Total Games: %d", 10))
	s.win.ColorOff(types.WHITE_BLACK)

	return nil
}

func (s *ScoreBoard) drawScoreBoard() error {
	s.win.ColorOn(types.WHITE_BLACK)
	s.win.MovePrint(3, 2, fmt.Sprintf("Player: %d", s.CurrentScore.PlayerScore))
	s.win.MovePrint(4, 2, fmt.Sprintf("Enemy: %d", s.CurrentScore.EnemyScore))
	s.win.ColorOff(types.WHITE_BLACK)

	return nil
}

func (s *ScoreBoard) drawAliveShips() error {
	s.win.ColorOn(types.WHITE_BLACK)
	// s.win.MovePrint(5, 2, fmt.Sprintf("Player: %d", s.CurrentScore.PlayerScore))
	// s.win.MovePrint(6, 2, fmt.Sprintf("Enemy: %d", s.CurrentScore.EnemyScore))
	s.win.ColorOff(types.WHITE_BLACK)

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

	s.win.ColorOn(types.WHITE_BLACK)
	s.win.MovePrint(startY, startX, title)
	s.win.ColorOff(types.WHITE_BLACK)

	// How do i render the borders properly?
	// Why is this a problem? Cuz rows and columns have variable length content
	for col := 0; col < cols; col++ {
		for row := 0; row < rows; row++ {

			x := (startX) + col*(len(stats[row][col])+1) + 1
			y := (startY) + row + 1

			s.win.ColorOn(types.WHITE_BLACK)
			s.win.MovePrint(y, x, stats[row][col])
			s.win.ColorOff(types.WHITE_BLACK)
		}
	}

	return nil
}
