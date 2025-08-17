package main

import (
	"context"

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

	gs.PlayerMap = NewMap(calculateSubWindow(stdscr, types.PLAYER), nil)
	gs.EnemyMap = NewMap(calculateSubWindow(stdscr, types.ENEMY), nil)
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
			gs.PlayerMap.Close()
			gs.EnemyMap.Close()
			gs.ScoreBoard.Close()
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

			gs.EnemyMap.SetEntity(types.Entity{
				Type:          types.CRUISER,
				StartPosition: types.Position{X: 2, Y: 2},
				EndPosition:   types.Position{X: 6, Y: 6},
				Color:         types.COLOR_SHIP,
				Sprite: map[types.Orientation][]rune{
					types.HORIZONTAL: {'%', '%', '%', '%'},
					types.VERTICAL:   {'%', '%', '%', '%'},
				},
			},
				types.VERTICAL)
			gs.EnemyMap.Render(ctx)
			gs.PlayerMap.Render(ctx)
			gs.ScoreBoard.Render(ctx)

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
