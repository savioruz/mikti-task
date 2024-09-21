package cli

import (
	"bufio"
	"fmt"
	"github.com/savioruz/mikti-task-1/pkg/constants"
	"golang.org/x/term"
	"os"
	"strings"
)

type Box struct {
	TopLeft, TopRight, BottomLeft, BottomRight, Vertical, Horizontal, LeftT, RightT string
}

var Style = Box{
	TopLeft:     constants.TopLeft,
	TopRight:    constants.TopRight,
	BottomLeft:  constants.BottomLeft,
	BottomRight: constants.BottomRight,
	Vertical:    constants.VerticalBar,
	Horizontal:  constants.HorizontalBar,
	LeftT:       constants.LeftT,
	RightT:      constants.RightT,
}

func (c *CLI) GetTerminalWidth() int {
	width, err, _ := term.GetSize(0)
	if err != 0 {
		return 80
	}
	return width
}

func (c *CLI) DrawLine(Style Box, left, right string, width int, text string) {
	content := strings.Repeat(Style.Horizontal, width)
	if text != "" {
		textStart := (width - len(text) - 2) / 2
		content = fmt.Sprintf("%s %s %s", content[:textStart], text, content[textStart+len(text)+2:])
	}
	fmt.Printf("%s%s%s\n", left, content, right)
}

func (c *CLI) DrawBox(items []string, title string) {
	width := c.GetTerminalWidth()

	c.DrawLine(Style, Style.TopLeft, Style.TopRight, width-2, "")
	fmt.Printf("%s %-*s %s\n", Style.Vertical, width-4, title, Style.Vertical)
	c.DrawLine(Style, Style.LeftT, Style.RightT, width-2, "")

	for _, item := range items {
		fmt.Printf("%s %-*s %s\n", Style.Vertical, width-4, item, Style.Vertical)
	}

	c.DrawLine(Style, Style.BottomLeft, Style.BottomRight, width-2, "")
}

func (c *CLI) GetUserInput(prompt string) string {
	fmt.Printf("%2s%s\n%2s>%2s", " ", prompt, " ", " ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func (c *CLI) GetOutputWithoutLines(prompt string) {
	fmt.Printf("%2s%s", " ", prompt)
}

func (c *CLI) GetOutput(prompt string) {
	fmt.Println()
	c.DrawLine(Style, Style.TopLeft, Style.TopRight, len(prompt)+2, "")
	fmt.Printf("%s %s %s\n", Style.Vertical, prompt, Style.Vertical)
	c.DrawLine(Style, Style.BottomLeft, Style.BottomRight, len(prompt)+2, "")
}

func (c *CLI) ClearScreen() {
	fmt.Print("\033[H\033[2J")
}

func (c *CLI) OneLineUp() {
	fmt.Printf("\033[1A")
}
