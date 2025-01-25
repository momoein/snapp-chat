package cli

import (
	app "snappchat/app/client"

	"github.com/awesome-gocui/gocui"
)

func keybindings(g *gocui.Gui, app app.App) error {
	var err error
	var (
		viewArr            = []string{ViewMenu, ViewMessageBar, ViewSendButton}
		cursorRequirements = []string{ViewMessageBar}
		active             = 0
	)

	// general handlers
	if err = g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err = g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, getNextView(viewArr, &active, cursorRequirements)); err != nil {
		return err
	}

	// send message button handler
	if err = g.SetKeybinding(ViewSendButton, gocui.KeyEnter, gocui.ModNone, getSendButtonHandler(app)); err != nil {
		return err
	}
	
	// menu handlers
	if err = g.SetKeybinding(ViewMenu, gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err = g.SetKeybinding(ViewMenu, gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err = g.SetKeybinding(ViewMenu, gocui.KeyEnter, gocui.ModNone, getMenuHandler(app)); err != nil { 
		return err
	}
	
	// window handler
	if err = g.SetKeybinding(ViewWindow, gocui.KeyEnter, gocui.ModNone, closeWindow); err != nil {
		return err
	}

	return nil
}
