package dotenv_test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/kick-project/maker/internal/resources/errs"
	"github.com/kick-project/maker/internal/resources/exit"
	"github.com/kick-project/maker/internal/resources/logger"
	"github.com/kick-project/maker/internal/resources/testtools"
	"github.com/kick-project/maker/internal/services/dotenv"
	"github.com/stretchr/testify/assert"
)

func makeDotenv() *dotenv.Dotenv {
	exitHandler := &exit.Handler{
		Mode: exit.MPanic,
	}
	logHandler := logger.New(os.Stderr, "ERROR", log.Ldate|log.Ltime|log.Lshortfile|log.Lmsgprefix, logger.InfoLevel, exitHandler)
	errHandler := errs.New(exitHandler, logHandler)

	d := dotenv.Defaults(
		&dotenv.Dotenv{
			Errs: errHandler,
			Exit: exitHandler,
		},
	)

	return d
}

func pathEnvfile() string {
	return filepath.Join(testtools.FixtureDir(), "menu", ".env")
}

func pathMakeFile() string {
	return filepath.Join(testtools.FixtureDir(), "menu", "Makefile")
}

func TestDotenv_Load(t *testing.T) {
	d := makeDotenv()
	os.Setenv("MAKER_GLOBAL", "From global")
	env1 := filepath.Join(testtools.FixtureDir(), "envfiles", ".env1")
	env2 := filepath.Join(testtools.FixtureDir(), "envfiles", ".env2")
	varfiles := fmt.Sprintf("%s,%s", env1, env2)
	d.Load(varfiles)

	assert.Equal(t, os.Getenv("MAKER_GLOBAL"), "From global")
	assert.Equal(t, os.Getenv("MAKER_VAR1"), "From .env1")
	assert.Equal(t, os.Getenv("MAKER_VAR2"), "From .env2")
}

func TestDotenv_WrapTarget(t *testing.T) {
	menuDir := filepath.Join(testtools.FixtureDir(), "menu")
	err := os.Chdir(menuDir)
	if err != nil {
		t.Error(err)
	}
	d := makeDotenv()
	assert.NotPanics(t, func() {
		d.WrapTarget(pathEnvfile(), pathMakeFile(), "target1")
	})
}

func TestDotenv_WrapTarget_NoTarget(t *testing.T) {
	menuDir := filepath.Join(testtools.FixtureDir(), "menu")
	err := os.Chdir(menuDir)
	if err != nil {
		t.Error(err)
	}
	d := makeDotenv()
	assert.Panics(t, func() {
		d.WrapTarget(pathEnvfile(), "Makefile", "notarget")
	})
}

func TestDotenv_WrapTarget_NoMakefile(t *testing.T) {
	noMakefile := filepath.Join(testtools.FixtureDir(), "missingmake")
	err := os.Chdir(noMakefile)
	if err != nil {
		t.Error(err)
	}
	d := makeDotenv()
	assert.Panics(t, func() {
		d.WrapTarget(pathEnvfile(), "Makefile", "notarget")
	})
}
