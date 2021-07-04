package menu

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"

	"github.com/dexterp/maker/internal/resources/errs"
)

// Menu provides menu options for make
type Menu struct {
	Stdout io.Writer
	Errs   errs.HandlerIface
}

type section struct {
	help    string
	targets []*target
	longest int // Longest target length
}

type target struct {
	title string
	help  string
}

// Display print make menu help
func (m *Menu) Display(path string) {
	reSection := regexp.MustCompile(`^###\s+([^#]+)$`)
	reTarget := regexp.MustCompile(`^_?(\S+?):.*?##\s+(.*)$`)

	file, err := os.Open(path)
	m.Errs.FatalF(`error opening file %s: %w`, err)

	curSection := &section{}
	sections := []*section{curSection}

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		txt := scanner.Text()
		sectionText := reSection.FindStringSubmatch(txt)
		if len(sectionText) == 2 {
			curSection = &section{
				help: sectionText[1],
			}
			sections = append(sections, curSection)
			continue
		}

		targetTxt := reTarget.FindStringSubmatch(txt)
		if len(targetTxt) == 3 {
			title := targetTxt[1]
			help := targetTxt[2]
			t := &target{
				title: title,
				help:  help,
			}
			l := len(title)
			if curSection.longest < l {
				curSection.longest = l
			}
			curSection.targets = append(curSection.targets, t)
			continue
		}
	}

	err = scanner.Err()
	m.Errs.FatalF(`error reading file %s: %w`, path, err)

	fmt.Fprint(m.Stdout, "Usage: make <target>\n\n")

	for _, sec := range sections {
		if sec.help != "" {
			fmt.Fprintf(m.Stdout, "### %s\n", sec.help)
		}
		for _, t := range sec.targets {
			fmt.Fprintf(m.Stdout, `  %-`+strconv.Itoa(sec.longest)+"s - %s\n", t.title, t.help)
		}
	}
}
