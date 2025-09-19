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
	win          *goncurses.Window
	keyInputChan chan goncurses.Key

	PlayerMap *Map
	EnemyMap  *Map

	ScoreBoard *ScoreBoard
	Guide      *Guide
	menuWindow *goncurses.Window

	playerHasSetShips bool
}

func NewGameState(stdscr *goncurses.Window, keyInputChan chan goncurses.Key) *GameState {
	gs := &GameState{
		win:          stdscr,
		keyInputChan: keyInputChan,
	}

	gs.PlayerMap = NewMap(
		calculateSubWindow(stdscr, types.PLAYER), // window
		true,                                     // isPlayerMap
		"PLAYER",                                 // title
		types.GREEN_BLACK,                        // titleColor
		nil,                                      // startingGrid
		nil,                                      // gridWidth
		nil,                                      // gridHeight
		true,                                     // enableKeyboard
	)
	gs.EnemyMap = NewMap(
		calculateSubWindow(stdscr, types.ENEMY), // window
		false,                                   // isPlayerMap
		"ENEMY",                                 // title
		types.RED_BLACK,                         // titleColor
		nil,                                     // startingGrid
		nil,                                     // gridWidth
		nil,                                     // gridHeight
		true,                                    // enableKeyboard
	)
	gs.ScoreBoard = NewScoreBoard(calculateSubWindow(stdscr, types.SCORE))
	gs.Guide = NewGuide(calculateSubWindow(stdscr, types.GUIDE))
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
		case key := <-gs.keyInputChan:
			gs.PlayerMap.HandleKeyInput(key)
			gs.EnemyMap.HandleKeyInput(key)

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

			if !gs.playerHasSetShips {
				gs.PlayerMap.EnableCursor(true)
				gs.EnemyMap.EnableCursor(false)
			} else {
				gs.PlayerMap.EnableCursor(false)
				gs.EnemyMap.EnableCursor(true)
			}

			err := gs.EnemyMap.Render(ctx)
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
			err = gs.Guide.Render(ctx)
			if err != nil {
				return err
			}

			goncurses.Update()
		}
	}
}

func (gs *GameState) CloseResources() error {
	var errs []error

	for _, closer := range []io.Closer{
		gs.PlayerMap,
		gs.EnemyMap,
		gs.ScoreBoard,
		gs.Guide,
	} {
		if err := closer.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	close(gs.keyInputChan)

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

	case types.GUIDE:
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
