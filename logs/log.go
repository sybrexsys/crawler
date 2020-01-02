package logs

import (
	"flag"
	stdlog "log"
	"os"

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"

	"github.com/sybrexsys/crawler/config"
)

var logFileName string

type Logger struct {
	*logrus.Logger
	needClose bool
	handler   *os.File
}

func NewLog(cfg *config.Config) (*Logger, error) {
	log := &Logger{
		Logger: logrus.New(),
	}
	level, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		level = logrus.TraceLevel
	}
	log.SetLevel(level)
	log.SetFormatter(&prefixed.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.9999",
	})
	fn := logFileName
	if fn == "" {
		fn = cfg.LogfileName
	}
	if fn != "" {
		f, err := os.OpenFile(fn, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		f.Write([]byte{13, 13, 13, 13})
		log.SetOutput(f)
		stdlog.SetOutput(log.Writer())
		log.needClose = true
		log.handler = f
		return log, nil
	}
	log.SetOutput(os.Stdout)
	stdlog.SetOutput(log.Writer())
	return log, nil
}

func (log *Logger) Close() {
	if log.needClose {
		log.handler.Close()
	}
}

func init() {
	flag.StringVar(&logFileName, "l", "", "log filename")
}
