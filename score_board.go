package main

import (
	"context"

	"github.com/rthornton128/goncurses"
)

type ScoreBoard struct {
	win          *goncurses.Window
	CurrentScore Score
	OverallScore Score
}

type Score struct {
	PlayerScore int
	EnemyScore  int
}

func NewScoreBoard(win *goncurses.Window) *ScoreBoard {
	return &ScoreBoard{
		win: win,
	}
}

func (s *ScoreBoard) Render(ctx context.Context) error {

	s.win.Erase()
	s.win.Box(goncurses.ACS_VLINE, goncurses.ACS_HLINE)
	s.win.NoutRefresh()
	s.draw()

	return nil

}

func (s *ScoreBoard) Close() error {
	return s.win.Delete()
}

func (s *ScoreBoard) draw() error {
	return nil
}
