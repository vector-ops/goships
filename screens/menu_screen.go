package screens

import (
	"context"
	"log"

	"github.com/rthornton128/goncurses"
	"github.com/vector-ops/goships/types"
)

func ShowMenuScreen(ctx context.Context, menuwin *goncurses.Window) types.GameType {

	menu_items := map[types.GameType]string{types.PVE: "play against a bot", types.PVP: "play against other players", types.QUIT: "quit"}
	items := make([]*goncurses.MenuItem, len(menu_items))

	i := 0
	var err error
	for name, desc := range menu_items {
		items[i], err = goncurses.NewItem(string(name), desc)
		if err != nil {
			log.Fatal(err)
		}
		defer items[i].Free()
		i++
	}

	menuwin.Keypad(true)

	menu, err := goncurses.NewMenu(items)
	if err != nil {
		menuwin.Print(err)
		return types.QUIT
	}
	defer menu.Free()

	menu.SetWindow(menuwin)
	dwin := menuwin.Derived(6, 38, 3, 1)
	menu.SubWindow(dwin)
	menu.Mark(" * ")

	// Print centered menu title
	_, x := menuwin.MaxYX()
	title := "Main Menu"
	menuwin.Box(0, 0)
	// menuwin.ColorOn(1)
	menuwin.MovePrint(1, (x/2)-(len(title)/2), title)
	// menuwin.ColorOff(1)
	menuwin.MoveAddChar(2, 0, goncurses.ACS_LTEE)
	menuwin.HLine(2, 1, goncurses.ACS_HLINE, x-2)
	menuwin.MoveAddChar(2, x-1, goncurses.ACS_RTEE)

	menu.Post()
	defer menu.UnPost()
	menuwin.Refresh()

	for {
		goncurses.Update()

		ch := menuwin.GetChar()

		switch goncurses.KeyString(ch) {
		case "q":
			return types.QUIT

		case "down":
			menu.Driver(goncurses.REQ_DOWN)
		case "up":
			menu.Driver(goncurses.REQ_UP)
		case "enter", "return":
			current := menu.Current(nil)
			name := current.Name()
			for k := range menu_items {
				if string(k) == name {
					return k
				}
			}

		}
	}

}
