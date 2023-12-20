package main

import (
	"os"
	"syscall"

	"github.com/gdamore/tcell/v2"
	"github.com/jessehorne/gostalgia/engine"
)

func main() {
	s, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	if err := s.Init(); err != nil {
		panic(err)
	}

	engine.Init(s)
	for {
		switch ev := s.PollEvent().(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape {
				if engine.Cmd != nil {
					engine.Cmd.Process.Signal(os.Interrupt)
					engine.Cmd.Process.Signal(syscall.SIGINT)
					engine.Cmd.Process.Signal(syscall.SIGTERM)
					engine.Cmd.Process.Kill()
					engine.Cmd = nil
				} else {
					s.Fini()
					os.Exit(0)
				}
			} else if ev.Key() == tcell.KeyBackspace {
				engine.DoBackspace(s)
			} else if ev.Key() == tcell.KeyEnter {
				engine.DoEnter(s)
			} else {
				r := ev.Rune()
				engine.HandleInput(s, r)
			}
		}
	}
}
