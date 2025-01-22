package cli

import (
	"errors"
	"fmt"

	"github.com/awesome-gocui/gocui"
)

type Widget interface {
	Layout(*gocui.Gui) error
}

type ButtonWidget struct {
	name    string
	title   string
	label   string
	x, y    int
	w, h    int
	handler func(g *gocui.Gui, v *gocui.View) error
}

func NewButtonWidget(name, title, label string, x, y, w, h int, handler func(g *gocui.Gui, v *gocui.View) error) Widget {
	return &ButtonWidget{
		name:    name,
		title:   title,
		label:   label,
		x:       x,
		y:       y,
		w:       w,
		h:       h,
		handler: handler,
	}
}

func (w *ButtonWidget) Layout(g *gocui.Gui) error {
	v, err := g.SetView(w.name, w.x, w.y, w.x+w.w, w.y+w.h, 0)
	if err != nil && !errors.Is(err, gocui.ErrUnknownView) {
		return err
	}

	v.Title = w.title
	v.BgColor = gocui.ColorBlack
	v.FgColor = gocui.ColorGreen

	if err := g.SetKeybinding(w.name, gocui.KeyEnter, gocui.ModNone, w.handler); err != nil {
		return err
	}
	fmt.Fprint(v, w.label)
	return nil
}

type TextWidget struct {
	Name       string
	Title      string
	Overlaps   byte
	Editable   bool
	Wrap       bool
	Autoscroll bool
	X, Y       int
	W, H       int
	Handler    Handler
}

func NewTextWidget(name, title string, editable bool, x, y, w, h int, handler Handler) Widget {
	return &TextWidget{
		Name:       name,
		Title:      title,
		Overlaps:   0,
		Editable:   editable,
		Wrap:       false,
		Autoscroll: false,
		X:          x,
		Y:          y,
		W:          w,
		H:          h,
		Handler:    handler,
	}
}

func (w *TextWidget) Layout(g *gocui.Gui) error {
	v, err := g.SetView(w.Name, w.X, w.Y, w.W, w.H, w.Overlaps)

	if err != nil && !errors.Is(err, gocui.ErrUnknownView) {
		return fmt.Errorf("error on set view: %v", err)
	}

	v.Title = w.Title
	v.Editable = w.Editable
	v.Wrap = true
	v.Autoscroll = true

	return nil
}
