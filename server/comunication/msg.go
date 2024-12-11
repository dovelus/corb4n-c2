package comunication

import (
	"errors"
	"os"
	"time"

	"github.com/charmbracelet/log"
)

var Logger = log.NewWithOptions(os.Stderr, log.Options{
	ReportCaller:    true,
	ReportTimestamp: true,
	TimeFormat:      time.RFC3339,
})

var ErrNoResults error = errors.New("no results found")
var ErrImplantExists error = errors.New("implant already exists")
