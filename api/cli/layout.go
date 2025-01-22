package cli

import (
	"errors"
	"fmt"

	"github.com/awesome-gocui/gocui"
)

const (
	ViewMenu       = "menu"
	ViewChatRoom   = "chatRoom"
	ViewMessageBar = "messageBar"
	ViewSendButton = "sendButton"
)


func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView(ViewMenu, 0, 0, maxX/4-1, maxY-1, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = "menu"
		if _, err = setCurrentViewOnTop(g, ViewMenu); err != nil {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		fmt.Fprintln(v, "Item 1")
		fmt.Fprintln(v, "Item 2")
		fmt.Fprint(v, "Item 3\nItem 4")
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
