package engine

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
)

var validKeys = " abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`1234567890-=[]\\;',./~!@#$%^&*()_+{}|:\"<>?"
var defaultStyle = tcell.StyleDefault.Foreground(tcell.ColorGreen).Background(tcell.ColorBlack)
var currentLine = ""

func Init(s tcell.Screen) {
	s.SetStyle(defaultStyle)

	_, h := s.Size()

	Print(s, 1, 2, "Welcome to Gostalgia.")
	Print(s, 1, 3, "https://github.com/jessehorne/gostalgia")
	Print(s, 1, 4, "Try 'HELP', 'CLEAR', 'LIST', 'EXIT', 'NEW' or 'RUN'...")
	Print(s, 1, 5, "If you run an infinite loop, hit ESC to end it. :-)")
	Print(s, 0, h-1, "] ")
}

func Print(s tcell.Screen, x, y int, str string) {
	for _, c := range str {
		s.SetContent(x, y, c, nil, defaultStyle)
		x++
	}

	s.Show()
}

func HandleInput(s tcell.Screen, k rune) {
	w, _ := s.Size()

	// check if key is in valid typeable keys
	for _, validKey := range validKeys {
		if k == validKey {
			if len(currentLine) < w-2 {
				currentLine += string(k)
				break
			}
		}
	}

	RedrawInput(s)
}

func RedrawInput(s tcell.Screen) {
	w, h := s.Size()

	// clear input line
	for i := 0; i < w; i++ {
		Print(s, i, h-1, " ")
	}

	Print(s, 0, h-1, "] "+currentLine)
}

func DoBackspace(s tcell.Screen) {
	if len(currentLine) == 0 {
		return
	}

	currentLine = currentLine[:len(currentLine)-1]
	RedrawInput(s)
}

var lines = []string{}

func PrintOne(s tcell.Screen, l string) {
	// add lines to lines buffer
	lines = append(lines, l)

	_, h := s.Size()

	// redraw all lines
	s.Clear()

	currentY := h - 2
	for i := len(lines) - 1; i >= 0; i-- {
		Print(s, 0, currentY, lines[i])
		currentY--
		if currentY < (0) {
			break
		}
	}

	currentLine = ""

	RedrawInput(s)
}

func PrintLines(s tcell.Screen, l []string) {
	// add lines to lines buffer
	lines = append(lines, "] "+currentLine)
	lines = append(lines, l...)

	_, h := s.Size()

	// redraw all lines
	s.Clear()

	currentY := h - 2
	for i := len(lines) - 1; i >= 0; i-- {
		Print(s, 0, currentY, lines[i])
		currentY--
		if currentY < (0) {
			break
		}
	}

	currentLine = ""

	RedrawInput(s)
}

func DoEnter(s tcell.Screen) {
	if currentLine == "HELP" {
		CommandHelp(s)
	} else if currentLine == "CLEAR" {
		CommandClear(s)
	} else if currentLine == "EXIT" {
		CommandExit(s)
	} else if currentLine == "LIST" {
		CommandList(s)
	} else if currentLine == "RUN" {
		CommandRun(s)
	} else if currentLine == "NEW" {
		CommandNew(s)
	} else {
		// add code to list
		CommandAdd(s)
	}
}

func CommandHelp(s tcell.Screen) {
	PrintLines(s, []string{
		"Welcome to Gostalgia!",
		"",
		"The valid commands are 'HELP', 'CLEAR', 'LIST', 'EXIT', 'NEW' and 'RUN'...",
		"",
		"CLEAR = clear screen",
		"LIST = list the current program",
		"RUN = run the current program",
		"",
		"Also, feel free to type any valid Go statements.",
		"Don't forget the line numbers!",
		"",
		"Example:",
		"",
		"10 fmt.Println(\"hello world\")",
		"RUN",
	})
}

func CommandClear(s tcell.Screen) {
	s.Clear()
	currentLine = ""
	lines = []string{}
	RedrawInput(s)
}

func CommandExit(s tcell.Screen) {
	s.Fini()
	os.Exit(0)
}

type CodeLine struct {
	Line int
	Code string
}

var code = []CodeLine{}

func SortCode() {
	sort.Slice(code, func(i, j int) bool {
		return code[i].Line < code[j].Line
	})
}

func CommandList(s tcell.Screen) {
	SortCode()

	var newLines []string

	for _, c := range code {
		lineNumber := strconv.Itoa(c.Line)
		newLines = append(newLines, fmt.Sprintf("%s %s", lineNumber, c.Code))
	}

	PrintLines(s, newLines)
}

func CommandRun(s tcell.Screen) {
	// put current code in template
	SortCode()
	var codeLines string
	for _, c := range code {
		codeLines += c.Code + "\n"
	}

	runCode := fmt.Sprintf(template, codeLines)

	execute(s, runCode)
}

func CommandAdd(s tcell.Screen) {
	// add current command to list of codez
	splitted := strings.Split(currentLine, " ")
	number := splitted[0]
	numberConv, err := strconv.Atoi(number)
	if err != nil {
		PrintLines(s, []string{
			"Invalid line number",
		})
		return
	}

	if len(splitted) == 1 {
		// if it's just a number, remove that line from code
		for i, c := range code {
			if c.Line == numberConv {
				code = append(code[:i], code[i+1:]...)
				break
			}
		}

		PrintLines(s, []string{})
		return
	}

	var found bool
	for i, c := range code {
		if c.Line == numberConv {
			code[i] = CodeLine{
				Line: numberConv,
				Code: strings.Join(splitted[1:], " "),
			}
			found = true
			break
		}
	}

	if !found {
		code = append(code, CodeLine{
			Line: numberConv,
			Code: strings.Join(splitted[1:], " "),
		})
	}

	PrintLines(s, []string{})
}

func CommandNew(s tcell.Screen) {
	code = []CodeLine{}
	PrintLines(s, []string{})
}
