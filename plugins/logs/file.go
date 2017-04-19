package plugin_logs

import (
	"io"
	"log"
	"os"

	"github.com/VirgilSecurity/virgild/coreapi"
	"github.com/namsral/flag"
	"github.com/pkg/errors"
)

var logfile string

func init() {
	flag.StringVar(&logfile, "filelogger-file", "-", "Path to log file ('-' - special parameter for colsole outut)")

	coreapi.RegisterLogger("file", makeFileLogger)
}

func makeFileLogger() (coreapi.Logger, error) {
	var (
		infoOut io.Writer = os.Stdout
		warnOut io.Writer = os.Stdout
		errOut  io.Writer = os.Stderr
	)
	if logfile != "-" {
		file, err := os.OpenFile(logfile, os.O_APPEND, 0666)
		if err != nil {
			return nil, errors.Wrapf(err, "Make file logger: open file (%v)", logfile)
		}
		infoOut = file
		warnOut = file
		errOut = file
	}

	return fileLogger{
		info: log.New(infoOut, "[INFO] ", log.LUTC|log.LstdFlags|log.Lmicroseconds),
		warn: log.New(warnOut, "[WARNING] ", log.LUTC|log.LstdFlags|log.Lmicroseconds),
		err:  log.New(errOut, "[ERROR] ", log.LUTC|log.LstdFlags|log.Lmicroseconds),
	}, nil
}

// FileLogger is simple logger
type fileLogger struct {
	info *log.Logger
	warn *log.Logger
	err  *log.Logger
}

func (l fileLogger) Info(format string, args ...interface{}) {
	l.info.Printf(format, args...)
}
func (l fileLogger) Warn(format string, args ...interface{}) {
	l.warn.Printf(format, args...)
}
func (l fileLogger) Err(format string, args ...interface{}) {
	l.err.Printf(format, args...)
}
