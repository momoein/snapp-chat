package cli

import (
	"errors"
	"fmt"
	"log"
	app "snappchat/app/client"
	"time"

	"github.com/awesome-gocui/gocui"
)

type Handler func(g *gocui.Gui, v *gocui.View) error

func quit(*gocui.Gui, *gocui.View) error {
	return gocui.ErrQuit
}

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

func getNextView(viewArr []string, activeIdx *int, cursorRequirements []string) Handler {
	return func(g *gocui.Gui, v *gocui.View) error {
		if len(viewArr) == 0 {
			return nil
		}

		nextIndex := (*activeIdx + 1) % len(viewArr)
		name := viewArr[nextIndex]

		_, err := setCurrentViewOnTop(g, name)
		if err != nil {
			return err
		}

		for _, viewName := range cursorRequirements {
			if viewName == name {
				g.Cursor = true
			} else {
				g.Cursor = false
			}
		}

		*activeIdx = nextIndex
		return nil
	}
}

func getSendButtonHandler(app app.App) Handler {
	return func(g *gocui.Gui, v *gocui.View) error {
		msgView, err := g.View(ViewMessageBar)
		if err != nil {
			return err
		}

		message := msgView.Buffer()
		msgView.Clear()

		if err := app.Service().SendMessageWS(message); err != nil {
			msgView.WriteString(message)
			log.Println("failed to send message: ", err)
		}

		return nil
	}
}

func UpdateChat(app app.App, g *gocui.Gui, name string) {
	ticker := time.NewTicker(100 * time.Millisecond)
	for {
		<-ticker.C

		msgByte, err := app.Service().ReadMessageWS()
		if err != nil || msgByte == nil {
			continue
		}

		g.Update(func(g *gocui.Gui) error {
			v, err := g.View(name)
			if err != nil {
				return fmt.Errorf("error on get view (%v): %v", name, err)
			}
			v.Autoscroll = true

			// update chat whit new message
			_, err = fmt.Fprintln(v, string(msgByte))
			if err != nil {
				return fmt.Errorf("error on update chat view: %v", err)
			}
			return nil
		})
	}

}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

func getMenuHandler(app app.App) Handler {
	return func(g *gocui.Gui, v *gocui.View) error {
		var opt string
		var err error

		_, cy := v.Cursor()
		if opt, err = v.Line(cy); err != nil {
			opt = ""
		}

		switch opt {
		case OptJoinRoom:
			app.Service().JoinRoom()
		case OptLeaveRoom:
			app.Service().LeaveRoom()
		case OptUsers:
			return showUsers(app, g)
		case OptExit:
			app.Service().LeaveRoom()
			return gocui.ErrQuit
		}

		return nil
	}
}

func showUsers(app app.App, g *gocui.Gui) error {
	var users string

	usersId, err := app.Service().GetUsers()
	if err != nil {
		users = "cant fetch online users!"
	}

	for _, id := range usersId {
		users = users + fmt.Sprintf("user-%d\n", id)
	}

	maxX, maxY := g.Size()
	if v, err := g.SetView(ViewUsers, maxX/2-30, maxY/2, maxX/2+30, maxY/2+5, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		g.Cursor = true
		v.Title = "users"
		v.Editable = true
		fmt.Fprintln(v, users)
		if _, err := g.SetCurrentView(ViewUsers); err != nil {
			return err
		}
	}
	return nil
}

func delUsersView(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView(ViewUsers); err != nil {
		return err
	}
	if _, err := g.SetCurrentView(ViewMenu); err != nil {
		return err
	}
	g.Cursor = false
	return nil
}
