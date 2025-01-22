package cli

import (
	"errors"
	"log"
	"snappchat/internal/client"

	"github.com/awesome-gocui/gocui"
)

func Run(app *client.ClientApp) {
	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		log.Fatalf("Failed to initialize GUI: %v", err)
	}
	defer g.Close()

	g.Highlight = true
	g.SelFgColor = gocui.ColorGreen

	g.SetManagerFunc(layout)

	if err := keybindings(g, app); err != nil {
		log.Fatalf("Failed to set keybindings: %v", err)
	}

	go UpdateChat(app, g, ViewChatRoom)

	if err := g.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		log.Fatalf("Main loop error: %v", err)
	}
}

func keybindings(g *gocui.Gui, app *client.ClientApp) error {
	var err error
	var (
		viewArr            = []string{ViewMenu, ViewMessageBar, ViewSendButton}
		cursorRequirements = []string{ViewMessageBar}
		active             = 0
	)

	if err = g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err = g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, getNextView(viewArr, &active, cursorRequirements)); err != nil {
		return err
	}
	if err = g.SetKeybinding(ViewSendButton, gocui.KeyEnter, gocui.ModNone, getSendButtonHandler(app)); err != nil {
		return err
	}
	// TODO: menu
	if err = g.SetKeybinding(ViewMenu, gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err = g.SetKeybinding(ViewMenu, gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err = g.SetKeybinding(ViewMenu, gocui.KeyEnter, gocui.ModNone, getLine); err != nil {
		return err
	}
	if err = g.SetKeybinding("msg", gocui.KeyEnter, gocui.ModNone, delMsg); err != nil {
		return err
	}

	return nil
}
