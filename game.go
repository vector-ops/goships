package main

import (
	"context"
	"fmt"
	"io"

	"github.com/rthornton128/goncurses"
	"github.com/vector-ops/goships/screens"
	"github.com/vector-ops/goships/types"
	"github.com/vector-ops/goships/utils"
)

type GameState struct {
	win *goncurses.Window

	PlayerMap *Map
	EnemyMap  *Map

	ScoreBoard   *ScoreBoard
	bufferWindow *goncurses.Window
	menuWindow   *goncurses.Window
}

func NewGameState(stdscr *goncurses.Window) *GameState {
	gs := &GameState{
		win: stdscr,
	}

	gs.PlayerMap = NewMap(calculateSubWindow(stdscr, types.PLAYER), nil, nil, nil)
	gs.EnemyMap = NewMap(calculateSubWindow(stdscr, types.ENEMY), nil, nil, nil)
	gs.ScoreBoard = NewScoreBoard(calculateSubWindow(stdscr, types.SCORE))
	gs.bufferWindow = calculateSubWindow(stdscr, types.BUFFER)
	gs.menuWindow = calculateSubWindow(stdscr, types.MENU)
	return gs
}

func (gs *GameState) Render(ctx context.Context, cancel context.CancelFunc) error {
	loading := false
	loadingContext, loadingCancel := context.WithCancel(ctx)

	go func() {
		utils.Delay(2000)
		loading = false
		loadingCancel()
	}()

	for {
		select {
		case <-ctx.Done():
			gs.CloseResources()
			loadingCancel()
			return nil
		default:
			if loading {
				screens.ShowLoadingScreen(loadingContext, gs.win, nil)
				gs.win.Erase()
				gs.win.Refresh()
				screens.ShowWelcomeScreen(ctx, gs.win)
				gs.win.Erase()
				gs.win.Refresh()
			}

			// gameType := screens.ShowMenuScreen(ctx, gs.menuWindow)
			// gs.menuWindow.Delete()
			// gs.win.Erase()
			// gs.win.Refresh()

			// switch gameType {
			// case types.QUIT:
			// 	cancel()
			// default:
			// }

			err := gs.EnemyMap.SetEntity(types.Entity{
				Type:          types.BATTLESHIP,
				CellType:      types.CELL_BATTLESHIP,
				StartPosition: types.Position{X: 4, Y: 3},
				EndPosition:   types.Position{X: 7, Y: 3},
				Color:         types.COLOR_SHIP,
				Sprite: map[types.Orientation][]rune{
					types.HORIZONTAL: types.BATTLESHIP_SPRITE,
					types.VERTICAL:   types.BATTLESHIP_SPRITE,
				},
			},
				types.HORIZONTAL)
			if err != nil {
				return err
			}

			err = gs.EnemyMap.Render(ctx)
			if err != nil {
				return err
			}
			err = gs.PlayerMap.Render(ctx)
			if err != nil {
				return err
			}
			err = gs.ScoreBoard.Render(ctx)
			if err != nil {
				return err
			}

			renderBuffer(gs.bufferWindow)
			goncurses.Update()
		}
	}
}

func renderBuffer(win *goncurses.Window) {
	win.Erase()
	win.Box(goncurses.ACS_VLINE, goncurses.ACS_HLINE)
	win.NoutRefresh()
}

func (gs *GameState) CloseResources() error {
	var errs []error

	for _, closer := range []io.Closer{
		gs.PlayerMap,
		gs.EnemyMap,
		gs.ScoreBoard,
	} {
		if err := closer.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to close resources: %v", errs)
	}
	return nil
}

func calculateSubWindow(win *goncurses.Window, wType types.WindowType) *goncurses.Window {
	my, mx := win.MaxYX()

	var h, w, y, x int

	switch wType {
	case types.PLAYER:
		h = my / 2
		w = mx / 2
		y = (my / 2) + 1
		x = mx - (mx * 3 / 4)

	case types.ENEMY:
		h = my / 2
		w = mx / 2
		y = 0
		x = mx - (mx * 3 / 4)

	case types.SCORE:
		h = my
		w = mx / 4
		y = 0
		x = 0

	case types.BUFFER:
		h = my
		w = mx / 4
		y = 0
		x = mx - (mx / 4)

	case types.MENU:
		h = my / 2
		w = mx / 2
		y = 0
		x = mx - (mx * 3 / 4)

	}

	mwin, err := goncurses.NewWindow(h, w, y, x)
	if err != nil {
		panic(err)
	}

	return mwin
}
