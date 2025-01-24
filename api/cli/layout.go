package cli

import (
	"errors"
	"fmt"

	"github.com/awesome-gocui/gocui"
)

const (
	ViewMenu       = "menu"
	ViewWindow     = "window"
	ViewChatRoom   = "chatRoom"
	ViewMessageBar = "messageBar"
	ViewSendButton = "sendButton"
	ViewUsers      = "users"
)

type MenuOption = string

const (
	OptJoinRoom  MenuOption = "join room"
	OptLeaveRoom MenuOption = "leave room"
	OptUsers     MenuOption = "users"
	OptHelp      MenuOption = "help"
	OptExit      MenuOption = "exit"
)

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView(ViewMenu, 0, 0, maxX/4-1, maxY-1, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = ViewMenu
		if _, err = setCurrentViewOnTop(g, ViewMenu); err != nil {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorCyan
		v.SelFgColor = gocui.ColorBlack
		fmt.Fprintln(v, OptJoinRoom)
		fmt.Fprintln(v, OptLeaveRoom)
		fmt.Fprintln(v, OptUsers)
		fmt.Fprintln(v, OptHelp)
		fmt.Fprintln(v, OptExit)
	}

	chatRoom := NewTextWidget(ViewChatRoom, "chat room", false, maxX/4, 0, maxX-1, maxY-maxY/4-1, nil)
	if err := chatRoom.Layout(g); err != nil {
		return err
	}

	messageBar := NewTextWidget(ViewMessageBar, "message", true, maxX/4, maxY-maxY/4, maxX-10, maxY-1, nil)
	if err := messageBar.Layout(g); err != nil {
		return err
	}

	sendButton := NewButtonWidget(ViewSendButton, "send", "📨", maxX-9, maxY-3, 7, 2, nil)
	if err := sendButton.Layout(g); err != nil {
		return err
	}

	return nil
}

const helpText = `
Tab: Move between panels.
Ctrl+C: Exit the application.

Menu panel:
	- Arrow Keys: Navigate between options.
	- Enter: select highlighted option.
Send button:
	- Enter: send the message.
Users window:
	- Enter: close the window.
Help window:
	- Enter: close the window.
`
