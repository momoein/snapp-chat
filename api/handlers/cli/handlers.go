package cli

import (
	"errors"
	"fmt"
	"log"
	"snappchat/internal/client"

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

func getSendButtonHandler(app *client.ClientApp) Handler {
	return func(g *gocui.Gui, v *gocui.View) error {
		msgView, err := g.View(ViewMessageBar)
		if err != nil {
			return err
		}

		message := msgView.Buffer()
		msgView.Clear()

		if err := app.SendMessageWS(message); err != nil {
			msgView.WriteString(message)
			log.Println("failed to send message: ", err)
		}
		
		return nil
	}
}

func UpdateChat(app *client.ClientApp, g *gocui.Gui, name string) {
	for {
		msgByte, err := app.ReadMessageWS()
		if err != nil {
			log.Println("socket read error: ", err)
			continue
		}

		if msgByte == nil {
			continue
		}

		g.Update(func(g *gocui.Gui) error {
			v, err := g.View(name)
			if err != nil {
				return fmt.Errorf("error on get view (%v): %v", name, err)
			}
			v.Autoscroll = true

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

func getLine(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	maxX, maxY := g.Size()
	if v, err := g.SetView("msg", maxX/2-30, maxY/2, maxX/2+30, maxY/2+2, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		fmt.Fprintln(v, l)
		if _, err := g.SetCurrentView("msg"); err != nil {
			return err
		}
	}
	return nil
}

func delMsg(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("msg"); err != nil {
		return err
	}
	if _, err := g.SetCurrentView(ViewMenu); err != nil {
		return err
	}
	return nil
}
