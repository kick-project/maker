package dotenv

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kick-project/maker/internal/resources/dfaults"
	"github.com/kick-project/maker/internal/resources/errs"
)

// Dotenv shell wrapper for Makefile
type Dotenv struct {
	Errs   errs.HandlerIface
	Prefix string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func Defaults(dotenv *Dotenv) *Dotenv {
	if dotenv == nil {
		dotenv = &Dotenv{}
	}
	if dotenv.Stdin == nil {
		dotenv.Stdin = os.Stdin
	}
	if dotenv.Stdout == nil {
		dotenv.Stdout = os.Stdout
	}
	if dotenv.Stderr == nil {
		dotenv.Stderr = os.Stderr
	}
	return dotenv
}

func (i *Dotenv) Exec(envFiles string, args ...string) {
	var (
		command *exec.Cmd
	)
	i.Load(envFiles)
	switch len(args) {
	case 0:
		return
	case 1:
		command = exec.Command(args[0])
	default:
		command = exec.Command(args[0], args[1:]...)
	}

	command.Stdin = i.Stdin
	command.Stdout = i.Stdout
	command.Stderr = i.Stderr
	err := command.Run()
	i.Errs.Fatal(err)
}

func (i *Dotenv) Load(envFiles string) {
	envs := expandArgs(strings.Split(envFiles, ",")...)
	if len(envs) > 0 {
		_ = godotenv.Load(envs...)
	}
}

func (i *Dotenv) WrapTarget(envFiles, makefile, target string, args ...string) {
	t := i.hasTarget(makefile, target)
	if t == "" {
		return
	}
	i.Exec(envFiles, `make`, t)
}

func (i *Dotenv) hasTarget(makefile, target string) string {
	prefix := dfaults.String(`_`, i.Prefix)
	re1 := regexp.MustCompile(fmt.Sprintf(`^(%s):`, target))
	re2 := regexp.MustCompile(fmt.Sprintf(`^(%s%s):`, prefix, target))
	file, err := os.Open(makefile)
	i.Errs.Fatal(err)
	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		text := scanner.Text()
		atoms1 := re1.FindStringSubmatch(text)
		if atoms1 != nil {
			return atoms1[1]
		}

		atoms2 := re2.FindStringSubmatch(text)
		if atoms2 != nil {
			return atoms2[1]
		}
	}
	return ""
}

func expandArgs(envs ...string) (expanded []string) {
	for _, path := range envs {
		if strings.HasPrefix(path, "~/") {
			usr, _ := user.Current()
			dir := usr.HomeDir
			path = filepath.Join(dir, path[2:])
		}
		if info, err := os.Stat(path); err == nil {
			mode := info.Mode()
			if mode.IsRegular() {
				expanded = append(expanded, path)
			}
		}
	}
	return
}
