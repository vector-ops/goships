package main

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/rthornton128/goncurses"
	"github.com/vector-ops/goships/logger"
	"github.com/vector-ops/goships/types"
	"github.com/vector-ops/goships/utils"
)

type GameState struct {
	win          *goncurses.Window
	debug        bool
	keyInputChan chan goncurses.Key
	logChan      chan logger.Log
	logger       *logger.Logger

	PlayerMap *Map
	EnemyMap  *Map

	ScoreBoard *ScoreBoard
	Guide      *Guide
	LogWindow  *LogWindow
	menuWindow *goncurses.Window

	playerHasSetShips bool
}

func NewGameState(stdscr *goncurses.Window, keyInputChan chan goncurses.Key, debug bool) *GameState {

	l, err := logger.NewLogger(logger.WithNoOp())
	if err != nil {
		panic(err)
	}

	gs := &GameState{
		win:          stdscr,
		debug:        debug,
		keyInputChan: keyInputChan,
		logger:       l,
	}

	if debug {
		logCh := make(chan logger.Log, 1)

		l, err := logger.NewLogger(logger.WithLogChan(logCh))
		if err != nil {
			panic(err)
		}

		gs.LogWindow = NewLogWindow(calculateSubWindow(stdscr, types.LogWindow, debug), logCh)
		gs.logger = l
		gs.logChan = logCh
	}

	gs.PlayerMap = NewMap(
		calculateSubWindow(stdscr, types.Player, debug), // window
		true,             // isPlayerMap
		"PLAYER",         // title
		types.GreenBlack, // titleColor
		nil,              // startingGrid
		nil,              // gridWidth
		nil,              // gridHeight
		true,             // enableKeyboard
		debug,
		gs.logger,
	)

	gs.EnemyMap = NewMap(
		calculateSubWindow(stdscr, types.Enemy, debug), // window
		false,          // isPlayerMap
		"ENEMY",        // title
		types.RedBlack, // titleColor
		nil,            // startingGrid
		nil,            // gridWidth
		nil,            // gridHeight
		true,           // enableKeyboard
		debug,
		gs.logger,
	)

	gs.ScoreBoard = NewScoreBoard(calculateSubWindow(stdscr, types.Score, debug), []StatBoard{
		{Title: "SCORE", StatHeader: []string{"Player", "Enemy"}, StatValues: []string{"0", "0"}},
		{Title: "PLAYER", StatHeader: []string{"Hits", "Misses"}, StatValues: []string{"0", "0"}},
		{Title: "ENEMY", StatHeader: []string{"Hits", "Misses"}, StatValues: []string{"0", "0"}},
	}, debug)

	gs.Guide = NewGuide(calculateSubWindow(stdscr, types.Guide, debug), debug)

	gs.menuWindow = calculateSubWindow(stdscr, types.Menu, debug)

	gs.EnemyMap.LogOn(ShipsPlaced)
	gs.EnemyMap.LogOn(Hit)
	gs.PlayerMap.LogOn(ShipsPlaced)
	gs.PlayerMap.LogOn(Hit)

	// l.Infof("The whale is a huge mammal living in the ocean.")
	// l.Infof("Sperm whales fight and eat giant squids.")
	gs.logger.Infof("Game started")

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

			if gs.debug {
				if err := gs.LogWindow.Render(ctx); err != nil {
					return err
				}
			}

			if err := gs.EnemyMap.Render(ctx); err != nil {
				return err
			}

			if err := gs.PlayerMap.Render(ctx); err != nil {
				return err
			}

			if err := gs.ScoreBoard.Render(ctx); err != nil {
				return err
			}

			if err := gs.Guide.Render(ctx); err != nil {
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

	if gs.LogWindow != nil {
		gs.LogWindow.Close()
	}

	close(gs.keyInputChan)

	if gs.logChan != nil {
		close(gs.logChan)
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to close resources: %v", errs)
	}

	return nil
}

func calculateSubWindow(win *goncurses.Window, wType types.WindowType, debug bool) *goncurses.Window {
	my, mx := win.MaxYX()

	var h, w, y, x int

	switch wType {
	case types.Player:
		h = my / 2
		w = mx / 2
		y = my / 2
		x = mx - (mx * 3 / 4)

	case types.Enemy:
		h = my / 2
		w = mx / 2
		y = 0
		x = mx - (mx * 3 / 4)

	case types.Score:
		h = my
		w = mx / 4
		y = 0
		x = 0

	case types.Guide:
		if debug {
			h = my / 2
		} else {
			h = my
		}
		w = mx / 4
		y = 0
		x = mx - (mx / 4)

	case types.Menu:
		h = my / 2
		w = mx / 2
		y = 0
		x = mx - (mx * 3 / 4)
	case types.LogWindow:
		h = my / 2
		w = mx / 4
		y = my / 2
		x = mx - (mx / 4)
	}

	mwin, err := goncurses.NewWindow(h, w, y, x)
	if err != nil {
		panic(err)
	}

	return mwin
}
