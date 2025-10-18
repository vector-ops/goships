package main

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/rthornton128/goncurses"
	"github.com/vector-ops/goships/types"
	"github.com/vector-ops/goships/utils"
)

type GameState struct {
	win          *goncurses.Window
	debug        bool
	keyInputChan chan goncurses.Key

	PlayerMap *Map
	EnemyMap  *Map

	ScoreBoard *ScoreBoard
	Guide      *Guide
	menuWindow *goncurses.Window

	playerHasSetShips bool
}

func NewGameState(stdscr *goncurses.Window, keyInputChan chan goncurses.Key, debug bool) *GameState {
	gs := &GameState{
		win:          stdscr,
		debug:        debug,
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
		debug,
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
		debug,
	)
	gs.ScoreBoard = NewScoreBoard(calculateSubWindow(stdscr, types.SCORE), map[string]*StatBoard{
		"SCORE":  {Title: "SCORE", StatHeader: []string{"Player", "Enemy"}, StatValues: []string{"0", "0"}},
		"PLAYER": {Title: "PLAYER", StatHeader: []string{"Hits", "Misses"}, StatValues: []string{"0", "0"}},
		"ENEMY":  {Title: "ENEMY", StatHeader: []string{"Hits", "Misses"}, StatValues: []string{"0", "0"}},
	}, debug)
	gs.Guide = NewGuide(calculateSubWindow(stdscr, types.GUIDE), debug)
	gs.menuWindow = calculateSubWindow(stdscr, types.MENU)
	return gs
}

func (gs *GameState) Render(ctx context.Context, cancel context.CancelFunc) error {
	for {
		select {
		case <-ctx.Done():
			gs.CloseResources()
			return nil
		case key := <-gs.keyInputChan:
			gs.PlayerMap.HandleKeyInput(key)
			gs.EnemyMap.HandleKeyInput(key)

		default:

			if !gs.playerHasSetShips {
				if gs.PlayerMap.HasPlacedShips() {
					gs.playerHasSetShips = true

					err := gs.EnemyMap.PlaceRandomShips()
					if err != nil {
						utils.WriteError(err)
					}
					gs.PlayerMap.EnableCursor(false)
					gs.EnemyMap.EnableCursor(true)
				} else {
					if gs.debug {
						gs.EnemyMap.SaveState()
					}
					gs.PlayerMap.EnableCursor(true)
					gs.EnemyMap.EnableCursor(false)
				}
			}

			enemyStats := gs.EnemyMap.GetStats()
			playerStats := gs.PlayerMap.GetStats()
			gs.ScoreBoard.SetStat("PLAYER", []string{strconv.Itoa(enemyStats.Hits), strconv.Itoa(enemyStats.Misses)})
			gs.ScoreBoard.SetStat("ENEMY", []string{strconv.Itoa(playerStats.Hits), strconv.Itoa(playerStats.Misses)})

			gs.ScoreBoard.SetStat("SCORE", []string{strconv.Itoa(enemyStats.GetShipsDestroyed()), strconv.Itoa(playerStats.GetShipsDestroyed())})

			if gs.playerHasSetShips {
				if gs.PlayerMap.GetTurn() < gs.EnemyMap.GetTurn() {
					gs.PlayerMap.HitRandomSpot()
				}
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
