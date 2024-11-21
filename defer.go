package structlog

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/mitchellh/panicwrap"
)

type deferredMsg struct {
	method  func(*Logger, any, ...any)
	msg     any
	keyvals []any
}

// DeferredLogger allows to defer logging until logger will be configured.
//
// Calls to logging methods (Debug, Info, Warn, PrintErr, Err, Print*, Fatal*) will enqueue
// message to log if called before Fatal* and (if WrapPanic was called) in child process.
// Logging methods called either after Fatal* or (if WrapPanic was called) in parent process
// will be ignored.
//
// Execute will output queued messages and call os.Exit if Fatal* was called or (if WrapPanic
// was called) in parent process.
//
// WARNING! Unlike usual loggers Fatal* methods won't call os.Exit, so execution of your code
// after Fatal method will be continued.
//
// Typical usage example:
//
//	func main() {
//	    var deferLog structlog.DeferredLogger
//	    deferLog.WrapPanic(&panicwrap.WrapConfig{HidePanic: true})
//
//	    err := someSetup()
//	    if err != nil {
//	        deferLog.Fatalf("failed to someSetup: %s", err) // Won't call os.Exit now!
//	    }
//
//	    var logJSON = flag.Bool("log.json", false, "use JSON log format")
//	    flag.Parse()
//	    if *logJSON {
//	        structlog.DefaultLogger.SetLogFormat(structlog.JSON)
//	    }
//
//	    log := structlog.New()
//	    deferLog.Execute(log) // os.Exit will happens here if someSetup has failed or code below will panic.
//	    // Do not use deferLog below this point.
//
//	    // ...
//	}
type DeferredLogger struct {
	mu          sync.Mutex
	exited      bool // True if we should had exited with exitStatus.
	exitStatus  int
	exitMsg     string // Message to be logged before exiting, if any.
	panicked    bool   // If true then exitMsg contains multiline panic output.
	msgs        []deferredMsg
	executed    bool
	userHandler func(string)
	exitFunc    func(int)
}

// SetExitFunc sets custom func to be called by Execute instead of os.Exit.
func (d *DeferredLogger) SetExitFunc(exit func(int)) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.exitFunc = exit
}

// WrapPanic uses github.com/mitchellh/panicwrap to catch and log panic.
// Logging will be deferred until Execute will be called.
//
// Set cfg.HidePanic=true to avoid duplicate output of panic to stderr.
func (d *DeferredLogger) WrapPanic(cfg *panicwrap.WrapConfig) {
	d.mu.Lock()
	if cfg == nil {
		cfg = &panicwrap.WrapConfig{}
	}
	d.userHandler = cfg.Handler
	cfg.Handler = d.handlePanic
	d.mu.Unlock()

	exitStatus, err := panicwrap.Wrap(cfg)
	d.mu.Lock()
	d.exitStatus = exitStatus
	d.mu.Unlock()

	switch {
	case err != nil:
		d.Fatalf("failed to panicwrap: %s", err)
	case !panicwrap.Wrapped(nil): // https://github.com/mitchellh/panicwrap/issues/18
		d.mu.Lock()
		d.exited = true
		d.mu.Unlock()
	}
}

func (d *DeferredLogger) handlePanic(output string) {
	d.mu.Lock()
	d.exitMsg = output
	d.panicked = true
	userHandler := d.userHandler
	d.mu.Unlock()

	if userHandler != nil {
		userHandler(output)
	}
}

// Execute must be called when log will be configured and ready to use. DeferredLogger must
// not be used after Execute.
//
// It will outputs deferred log messages followed by possible panic (if WrapPanic was called).
// If Fatal* was called or WrapPanic was called and panic happened calls os.Exit (or function
// set by SetExitFunc).
func (d *DeferredLogger) Execute(log *Logger) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.executed {
		panic("must be called just once")
	}
	d.executed = true

	for _, m := range d.msgs {
		m.method(log, m.msg, m.keyvals...)
	}

	if !d.exited {
		return
	}
	if d.exitMsg != "" {
		if d.panicked {
			log = log.New().SetKeysFormat(map[string]string{KeyMessage: " %[2]s"})
		}
		log.PrintErr(d.exitMsg)
	}

	if d.exitFunc == nil {
		d.exitFunc = os.Exit
	}
	d.exitFunc(d.exitStatus)
}

func (d *DeferredLogger) log(method func(*Logger, any, ...any), msg any, keyvals ...any) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.executed {
		panic("must not be called after Execute")
	}
	if !d.exited {
		d.msgs = append(d.msgs, deferredMsg{
			method:  method,
			msg:     msg,
			keyvals: keyvals,
		})
	}
}

func (d *DeferredLogger) fatal(msg string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.executed {
		panic("must not be called after Execute")
	}
	if !d.exited {
		d.exited = true
		d.exitStatus = 1
		d.exitMsg = msg
	}
}

// PrintErr defers logging until Execute will be called.
func (d *DeferredLogger) PrintErr(msg any, keyvals ...any) {
	d.log((*Logger).PrintErr, msg, keyvals...)
}

// Err defers logging until Execute will be called.
func (d *DeferredLogger) Err(msg any, keyvals ...any) error {
	d.log((*Logger).PrintErr, msg, keyvals...)
	return getErr(msg, keyvals...)
}

// Warn defers logging until Execute will be called.
func (d *DeferredLogger) Warn(msg any, keyvals ...any) {
	d.log((*Logger).Warn, msg, keyvals...)
}

// Info defers logging until Execute will be called.
func (d *DeferredLogger) Info(msg any, keyvals ...any) {
	d.log((*Logger).Info, msg, keyvals...)
}

// Debug defers logging until Execute will be called.
//
//nolint:godox // False positive.
func (d *DeferredLogger) Debug(msg any, keyvals ...any) {
	d.log((*Logger).Debug, msg, keyvals...)
}

// Print defers logging until Execute will be called.
func (d *DeferredLogger) Print(v ...any) {
	d.log((*Logger).Info, fmt.Sprint(v...))
}

// Printf defers logging until Execute will be called.
func (d *DeferredLogger) Printf(format string, v ...any) {
	d.log((*Logger).Info, fmt.Sprintf(format, v...))
}

// Println defers logging until Execute will be called.
func (d *DeferredLogger) Println(v ...any) {
	d.log((*Logger).Info, strings.TrimSuffix(fmt.Sprintln(v...), "\n"))
}

// Fatal defers logging until Execute will be called.
//
// WARNING: os.Exit won't be called until Execute, so execution of your code will continue.
func (d *DeferredLogger) Fatal(v ...any) {
	d.fatal(fmt.Sprint(v...))
}

// Fatalf defers logging until Execute will be called.
//
// WARNING: os.Exit won't be called until Execute, so execution of your code will continue.
func (d *DeferredLogger) Fatalf(format string, v ...any) {
	d.fatal(fmt.Sprintf(format, v...))
}

// Fatalln defers logging until Execute will be called.
//
// WARNING: os.Exit won't be called until Execute, so execution of your code will continue.
func (d *DeferredLogger) Fatalln(v ...any) {
	d.fatal(strings.TrimSuffix(fmt.Sprintln(v...), "\n"))
}
