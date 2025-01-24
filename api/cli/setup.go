package cli

import (
	"errors"
	"log"
	app "snappchat/app/client"

	"github.com/awesome-gocui/gocui"
)

func Run(app app.App) {
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
