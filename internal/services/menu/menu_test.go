package menu_test

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kick-project/maker/internal/di"
	"github.com/kick-project/maker/internal/resources/testtools"
	"github.com/stretchr/testify/assert"
)

func TestMenu_Display(t *testing.T) {
	fixtures := testtools.FixtureDir()
	makefile := filepath.Join(fixtures, "menu", "Makefile")

	buf := bytes.Buffer{}
	inject := di.Defaults(&di.DI{
		Stdout: &buf,
	})

	menu := inject.MakeMenu()
	menu.Display(makefile)

	txt := buf.String()
	lines := []string{
		`Usage: make <target>`,
		``,
		`### Section 1`,
		`  target1  - Target 1`,
		`  target22 - Target 2`,
		`### Section 2`,
		`  target333  - Target 3`,
		`  target4444 - Target 4`,
	}
	assert.Regexp(t, strings.Join(lines, `\n`), txt)
}
