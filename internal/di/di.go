package di

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/kick-project/maker/internal/resources/errs"
	"github.com/kick-project/maker/internal/resources/exit"
	"github.com/kick-project/maker/internal/resources/logger"
	"github.com/kick-project/maker/internal/services/dotenv"
	"github.com/kick-project/maker/internal/services/menu"
)

// DI dependency injection struct
type DI struct {
	Stderr   io.Writer    // Stderr
	Stdout   io.Writer    // Stdout
	LogLevel logger.Level // Logger level
	ExitMode int          // Exit mode one of exit.MNone or exit.Panic
	Prefix   string       // Target prefix

	/* Cache objects */
	cacheErrHandler  *errs.Handler
	cacheExitHandler *exit.Handler
	cacheLogFile     *os.File
	cacheMenu        *menu.Menu
	cacheDotenv      *dotenv.Dotenv
}

// Defaults add defaults to DI
func Defaults(di *DI) *DI {
	if di.Stderr == nil {
		di.Stderr = os.Stderr
	}
	if di.Stdout == nil {
		di.Stdout = os.Stdout
	}
	if di.Prefix == "" {
		di.Prefix = "_"
	}
	return di
}

//
// Dependency Injectors
//

// MakeShell injector
func (i *DI) MakeDotenv() *dotenv.Dotenv {
	if i.cacheDotenv != nil {
		return i.cacheDotenv
	}
	i.cacheDotenv = dotenv.Defaults(&dotenv.Dotenv{
		Errs:   i.MakeErrorHandler(),
		Prefix: i.Prefix,
	})
	return i.cacheDotenv
}

// MakeErrorHandler dependency injector
func (i *DI) MakeErrorHandler() *errs.Handler {
	if i.cacheErrHandler != nil {
		return i.cacheErrHandler
	}
	handler := errs.New(i.MakeExitHandler(), i.MakeLoggerOutput("", ""))
	i.cacheErrHandler = handler
	return handler
}

// MakeErrorHandler dependency injector
func (i *DI) MakeMenu() *menu.Menu {
	if i.cacheMenu != nil {
		return i.cacheMenu
	}
	i.cacheMenu = &menu.Menu{
		Stdout: i.Stdout,
		Errs:   i.MakeErrorHandler(),
	}
	return i.cacheMenu
}

// MakeExitHandler dependency injector
func (i *DI) MakeExitHandler() *exit.Handler {
	if i.cacheExitHandler != nil {
		return i.cacheExitHandler
	}
	handler := &exit.Handler{
		Mode: i.ExitMode,
	}
	i.cacheExitHandler = handler
	return handler
}

// MakeLogFile create a logfile and return the interface
func (i *DI) MakeLogFile(logfile string) *os.File {
	if i.cacheLogFile != nil {
		return i.cacheLogFile
	}
	var (
		f   *os.File
		err error
	)

	fInfo, err := os.Stat(logfile)
	if err != nil && !os.IsNotExist(err) {
		// Simple output because logging is not available
		fmt.Printf(`can not open log file %s: %v`, logfile, err)
	} else if err == nil && fInfo.Size() > 1024*1024*2 {
		// Remove files greater than 2M
		os.Remove(logfile)
	}
	f, err = os.OpenFile(logfile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		i.MakeErrorHandler().FatalF(`can not open logfile %s: %v`, logfile, err)
	}
	i.cacheLogFile = f
	return f
}

// MakeLoggerOutput inject logger.OutputIface.
func (i *DI) MakeLoggerOutput(prefix string, logfile string) *logger.Router {
	toStderr := logger.New(i.Stderr, prefix, log.Lmsgprefix, i.LogLevel, i.MakeExitHandler())
	if logfile != "" {
		toFile := logger.New(i.MakeLogFile(logfile), prefix, log.Ldate|log.Ltime|log.Lshortfile|log.Lmsgprefix, i.LogLevel, i.MakeExitHandler())
		return logger.NewRouter(toFile, toStderr)
	}
	return logger.NewRouter(toStderr)
}
