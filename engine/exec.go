package engine

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
)

var template = `
package main

import (
	"fmt"
)

func main() {
	%s
}
`

var Cmd *exec.Cmd

func execute(s tcell.Screen, code string) {
	PrintLines(s, []string{})
	
	// create temp file
	dname, err := os.MkdirTemp("", "code")
	//defer os.Remove(dname)
	if err != nil {
		PrintLines(s, []string{err.Error()})
		return
	}

	f, err := os.Create(dname + "/main.go")
	f.Chdir()
	if err != nil {
		return
	}
	defer f.Close()
	f.WriteString(code)

	init := exec.Command("go", "mod", "init", "temp")
	init.Dir = dname
	init.Run()
	
	Cmd = exec.Command(
		"go", "run", filepath.Base(f.Name()))
	
	Cmd.Dir = dname

	cmdReader, err := Cmd.StdoutPipe()
	if err != nil {
		PrintLines(s, []string{err.Error()})
		return
	}
	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			if Cmd == nil {
				break
			}
			PrintOne(s, scanner.Text())
		}
	}()

	errReader, err := Cmd.StderrPipe()
	if err != nil {
		PrintLines(s, []string{err.Error()})
		return
	}
	errscanner := bufio.NewScanner(errReader)
	go func() {
		for errscanner.Scan() {
			if Cmd == nil {
				break
			}
			PrintOne(s, errscanner.Text())
		}
	}()

	err = Cmd.Start()
	if err != nil {
		PrintLines(s, []string{err.Error()})
		return
	}
	
	//err = Cmd.Wait()
	//if err != nil {
	//	PrintLines(s, []string{err.Error()})
	//	return
	//}
}
