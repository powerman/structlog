// Package structlog provides structured logger which looks friendly as
// plain text (more like handcrafted vertical-aligned log lines which
// doesn't feels like key/value pairs) and as JSON (each important
// value in log line as separate key/value pair).
package structlog

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"
)

func init() {
	log.SetFlags(0)
}

type (
	logFormat byte
	logLevel  byte
)

const (
	Text logFormat = iota
	JSON
)

const (
	DBG logLevel = iota
	INF
	WRN
	ERR
)

const (
	DefaultLogFormat     = Text
	DefaultLogLevel      = DBG
	DefaultKeyValFormat  = ` %s=%v`
	DefaultTimeFormat    = time.StampMicro
	DefaultTimeValFormat = time.RFC3339Nano
	MissingValue         = "(MISSING)"
)

const (
	KeyTime    = "_t" // Key name used to output current time.
	KeyApp     = "_a" // Key name used to output app name.
	KeyPID     = "_p" // Key name used to output PID.
	KeyLevel   = "_l" // Key name used to output log level.
	KeyUnit    = "_u" // Key name used to output unit/module/package name.
	KeyMessage = "_m" // Key name used to output log message.
	KeyFunc    = "_f" // Key name used to output caller's function name.
	KeySource  = "_s" // Key name used to output caller's file and line.
	KeyStack   = "__" // Key name used to output multiline stack trace.
)

// Auto can be used as value for KeyUnit and KeyStack to automatically
// generate their values: caller package's directory name and full stack
// of the current goroutine.
const Auto = "\x00"

const unknown = "???"

// ParseLevel convert levelName from flag or config file into logLevel.
func ParseLevel(levelName string) logLevel {
	switch strings.ToLower(levelName) {
	case "err", "error", "fatal", "crit", "critical", "alert", "emerg", "emergency":
		return ERR
	case "wrn", "warn", "warning":
		return WRN
	case "inf", "info", "notice":
		return INF
	case "dbg", "debug", "trace":
		return DBG
	default:
		DefaultLogger.Err("failed", "levelName", levelName)
		return DBG
	}
}

func (l logLevel) String() string {
	switch l {
	case ERR:
		return "ERR"
	case WRN:
		return "WRN"
	case INF:
		return "inf"
	case DBG:
		return "dbg"
	default:
		return unknown
	}
}

func (l logLevel) MarshalJSON() ([]byte, error) {
	return []byte(`"` + l.String() + `"`), nil
}

type Logger struct {
	parent         *Logger
	format         *logFormat
	level          *logLevel
	keyValFormat   *string
	timeFormat     *string
	timeValFormat  *string
	callDepth      int
	defaultKeyvals map[string]interface{}
	prefixKeys     []string
	suffixKeys     []string
	keysFormat     map[string]string
	sync.RWMutex
}

var (
	// DefaultLogger provides sane defaults inherited by new logger
	// objects created with New(). Feel free to change it settings
	// when your app start.
	DefaultLogger = NewZeroLogger(
		KeyApp, path.Base(os.Args[0]),
		KeyPID, os.Getpid(),
	).SetPrefixKeys(
		KeyApp, KeyPID, KeyLevel, KeyUnit,
	).SetSuffixKeys(
		KeyFunc, KeySource, KeyStack,
	).SetKeysFormat(map[string]string{
		KeyApp:     "%[2]s",
		KeyPID:     "[%[2]d]",
		KeyLevel:   " %[2]s",
		KeyUnit:    " %[2]s:",
		KeyMessage: " %#[2]q",
		KeyFunc:    " \t@ %[2]s",
		KeySource:  "(%[2]s)",
		KeyStack:   "\n%[2]s",
	})
)

// NewZeroLogger creates and returns a new logger with empty settings.
func NewZeroLogger(defaultKeyvals ...interface{}) *Logger {
	var (
		format        = DefaultLogFormat
		level         = DefaultLogLevel
		keyValFormat  = DefaultKeyValFormat
		timeFormat    = DefaultTimeFormat
		timeValFormat = DefaultTimeValFormat
	)
	return (&Logger{
		parent:        nil,
		format:        &format,
		level:         &level,
		keyValFormat:  &keyValFormat,
		timeFormat:    &timeFormat,
		timeValFormat: &timeValFormat,
		callDepth:     2,
		defaultKeyvals: map[string]interface{}{
			KeyUnit:   Auto,    // must be non-nil to enable field
			KeyFunc:   unknown, // must be non-nil to enable field
			KeySource: unknown, // must be non-nil to enable field
		},
		prefixKeys: []string{},
		suffixKeys: []string{},
		keysFormat: map[string]string{},
	}).New(defaultKeyvals...)
}

// New creates and returns a new logger which inherits all settings from
// DefaultLogger.
func New(defaultKeyvals ...interface{}) *Logger {
	return DefaultLogger.New(defaultKeyvals...)
}

// New creates and returns a new logger which inherits all settings from l.
func (l *Logger) New(defaultKeyvals ...interface{}) *Logger {
	if l == nil {
		panic("New called on nil *Logger")
	}
	return (&Logger{
		parent:         l,
		callDepth:      0,
		defaultKeyvals: make(map[string]interface{}, 16),
		prefixKeys:     make([]string, 0, 16),
		suffixKeys:     make([]string, 0, 16),
		keysFormat:     make(map[string]string, 16),
	}).SetDefaultKeyvals(defaultKeyvals...)
}

// SetLogFormat changes log output format (default value is
// DefaultLogFormat).
//
// It doesn't creates a new logger, it returns l just for convenience.
func (l *Logger) SetLogFormat(format logFormat) *Logger {
	l.Lock()
	defer l.Unlock()
	l.format = &format
	return l
}

// SetLogLevel changes minimum required log level to output log
// (default value is DefaultLogLevel).
//
// It doesn't creates a new logger, it returns l just for convenience.
func (l *Logger) SetLogLevel(level logLevel) *Logger {
	l.Lock()
	defer l.Unlock()
	l.level = &level
	return l
}

// SetKeyValFormat changes fmt format string used to output key/value pair
// for keys which doesn't have custom format set by SetKeysFormat (default
// value is DefaultKeyValFormat).
//
// See SetKeysFormat for more details.
//
// It doesn't creates a new logger, it returns l just for convenience.
func (l *Logger) SetKeyValFormat(format string) *Logger {
	l.Lock()
	defer l.Unlock()
	l.keyValFormat = &format
	return l
}

// SetTimeFormat changes format for time.Time.Format used when output log
// time (default value is DefaultTimeFormat).
//
// It doesn't creates a new logger, it returns l just for convenience.
func (l *Logger) SetTimeFormat(format string) *Logger {
	l.Lock()
	defer l.Unlock()
	l.timeFormat = &format
	return l
}

// SetTimeValFormat changes format for time.Time.Format used when output
// time.Time values (default value is DefaultTimeValFormat).
//
// It doesn't creates a new logger, it returns l just for convenience.
func (l *Logger) SetTimeValFormat(format string) *Logger {
	l.Lock()
	defer l.Unlock()
	l.timeValFormat = &format
	return l
}

// AddCallDepth will add depth to amount of skipped stack frames while
// calculating default values for KeyUnit, KeyFunc and KeySource.
//
// Use it if you want to report from perspective of your caller.
//
// It doesn't creates a new logger, it returns l just for convenience.
func (l *Logger) AddCallDepth(depth int) *Logger {
	l.Lock()
	defer l.Unlock()
	l.callDepth += depth
	return l
}

// SetDefaultKeyvals add/replace values for keys in defaultKeyvals.
//
// The keyvals must be a list of key/value pairs, keys must be a string.
// In case of odd amount of elements in keyvals it'll log error and use
// MissingValue as value for last key. In case of non-string keys it'll
// log error and convert key to string.
//
// Keys in defaultKeyvals will provide default values for
// prefixKeys/suffixKeys, but these values will be used only if their
// key is included in prefixKeys or suffixKeys and same key won't be
// included within keyvals provided with log message.
//
// To delete keys from defaultKeyvals set their value to nil.
// This is very useful if unwanted key was inherited from parent logger.
//
// It doesn't creates a new logger, it returns l just for convenience.
func (l *Logger) SetDefaultKeyvals(keyvals ...interface{}) *Logger {
	if len(keyvals)%2 != 0 {
		l.New().AddCallDepth(getPackageDepth()).Err("odd keyvals")
		keyvals = append(keyvals, MissingValue)
	}
	for i := 0; i < len(keyvals); i += 2 {
		k, ok := keyvals[i].(string)
		if !ok {
			l.New().AddCallDepth(getPackageDepth()).SetKeyValFormat(" %#[2]v").Err("key is not string", "key", keyvals[i])
			k = fmt.Sprint(keyvals[i])
		}
		l.Lock()
		l.defaultKeyvals[k] = keyvals[i+1]
		l.Unlock()
	}
	return l
}

// SetPrefixKeys replace current prefixKeys for l.
//
// These keys will be output right after l's parent prefixKeys, if any.
//
// XXX Panics if will be called after using l (or logger created using
// l.New()) to log anything.
//
// It doesn't creates a new logger, it returns l just for convenience.
func (l *Logger) SetPrefixKeys(keys ...string) *Logger {
	l.Lock()
	defer l.Unlock()
	if l.parent == nil {
		panic("too late to reconfigure prefixKeys")
	}
	l.prefixKeys = make([]string, len(keys))
	copy(l.prefixKeys, keys)
	return l
}

// AppendPrefixKeys appends keys to current prefixKeys for l.
//
// XXX Panics if will be called after using l (or logger created using
// l.New()) to log anything.
//
// It doesn't creates a new logger, it returns l just for convenience.
func (l *Logger) AppendPrefixKeys(keys ...string) *Logger {
	l.Lock()
	defer l.Unlock()
	if l.parent == nil {
		panic("too late to reconfigure prefixKeys")
	}
	l.prefixKeys = append(l.prefixKeys, keys...)
	return l
}

// PrependSuffixKeys prepend keys to current suffixKeys for l.
//
// XXX Panics if will be called after using l (or logger created using
// l.New()) to log anything.
//
// It doesn't creates a new logger, it returns l just for convenience.
func (l *Logger) PrependSuffixKeys(keys ...string) *Logger {
	l.Lock()
	defer l.Unlock()
	if l.parent == nil {
		panic("too late to reconfigure suffixKeys")
	}
	l.suffixKeys = append(append([]string(nil), keys...), l.suffixKeys...)
	return l
}

// SetSuffixKeys replace current suffixKeys for l.
//
// These keys will be output just before l's parent suffixKeys, if any.
//
// XXX Panics if will be called after using l (or logger created using
// l.New()) to log anything.
//
// It doesn't creates a new logger, it returns l just for convenience.
func (l *Logger) SetSuffixKeys(keys ...string) *Logger {
	l.Lock()
	defer l.Unlock()
	if l.parent == nil {
		panic("too late to reconfigure suffixKeys")
	}
	l.suffixKeys = make([]string, len(keys))
	copy(l.suffixKeys, keys)
	return l
}

// SetKeysFormat add/replace custom fmt format string for keys.
// If key doesn't have custom format string then it will use format set
// using SetKeyValFormat (default value is DefaultKeyValFormat).
//
// These format strings will be used as fmt.Sprintf(format,key,val),
// so you can refer to key name and it value as %[1] and %[2] - this is
// very useful in case you wanna output only key value, without name.
//
// No extra spaces will be output between key/value pairs, so if you need
// some delimiters then include them inside format strings.
//
// It doesn't creates a new logger, it returns l just for convenience.
func (l *Logger) SetKeysFormat(keysFormat map[string]string) *Logger {
	l.Lock()
	defer l.Unlock()
	for k, v := range keysFormat {
		l.keysFormat[k] = v
	}
	return l
}

// IsInfo returns true if l's log level DBG or INF.
func (l *Logger) IsInfo() bool {
	l.RLock()
	defer l.RUnlock()
	if l.parent != nil {
		l.RUnlock()
		l.mergeParent()
		l.RLock()
	}
	return *l.level <= INF
}

// IsDebug returns true if l's log level DBG.
func (l *Logger) IsDebug() bool {
	l.RLock()
	defer l.RUnlock()
	if l.parent != nil {
		l.RUnlock()
		l.mergeParent()
		l.RLock()
	}
	return *l.level <= DBG
}

// Recover calls recover(), and if it returns non-nil, then log
// defaultKeyvals, value returned by recover() and keyvals with stack
// trace and level ERR plus stores value returned by recover() into err if
// err is not nil.
//
//   defer log.Recover(nil)
//   func PanicToErr() (err error) { defer log.Recover(&err); ... }
func (l *Logger) Recover(err *error, keyvals ...interface{}) {
	if e := recover(); e != nil {
		if err != nil {
			var ok bool
			if *err, ok = e.(error); !ok {
				*err = fmt.Errorf("%v", e)
			}
		}

		keyvals = append(keyvals, KeyStack, Auto)
		const panicDepth = 2 // runtime.call64, runtime.gopanic
		l.New().AddCallDepth(panicDepth).log(ERR, e, keyvals...)
	}
}

// PrintErr log defaultKeyvals, msg and keyvals with level ERR.
//
// In most cases you should use Err instead, to both log and handle error.
func (l *Logger) PrintErr(msg interface{}, keyvals ...interface{}) {
	l.log(ERR, msg, keyvals...)
}

// Err log defaultKeyvals, msg and keyvals with level ERR and returns
// first arg of error type or msg if there are no errors in args.
//
//   return log.Err("message to log", "error to log and return", err)
//   return log.Err(errors.New("error to log and return"), "error to log", err)
func (l *Logger) Err(msg interface{}, keyvals ...interface{}) error {
	l.log(ERR, msg, keyvals...)
	return getErr(msg, keyvals...)
}

// Warn log defaultKeyvals, msg and keyvals with level WRN.
func (l *Logger) Warn(msg interface{}, keyvals ...interface{}) {
	l.log(WRN, msg, keyvals...)
}

// Info log defaultKeyvals, msg and keyvals with level INF.
func (l *Logger) Info(msg interface{}, keyvals ...interface{}) {
	l.log(INF, msg, keyvals...)
}

// Debug log defaultKeyvals, msg and keyvals with level DBG.
func (l *Logger) Debug(msg interface{}, keyvals ...interface{}) {
	l.log(DBG, msg, keyvals...)
}

// Print works like log.Print. Use level INF.
// Also output defaultKeyvals for prefixKeys/suffixKeys.
func (l *Logger) Print(v ...interface{}) {
	l.log(INF, fmt.Sprint(v...))
}

// Printf works like log.Printf. Use level INF.
// Also output defaultKeyvals for prefixKeys/suffixKeys.
func (l *Logger) Printf(format string, v ...interface{}) {
	l.log(INF, fmt.Sprintf(format, v...))
}

// Println works like log.Println. Use level INF.
// Also output defaultKeyvals for prefixKeys/suffixKeys.
func (l *Logger) Println(v ...interface{}) {
	l.log(INF, strings.TrimSuffix(fmt.Sprintln(v...), "\n"))
}

// Fatal works like log.Fatal. Use level ERR.
// Also output defaultKeyvals for prefixKeys/suffixKeys.
func (l *Logger) Fatal(v ...interface{}) {
	l.log(ERR, fmt.Sprint(v...))
	os.Exit(1)
}

// Fatalf works like log.Fatalf. Use level ERR.
// Also output defaultKeyvals for prefixKeys/suffixKeys.
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.log(ERR, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// Fatalln works like log.Fatalln. Use level ERR.
// Also output defaultKeyvals for prefixKeys/suffixKeys.
func (l *Logger) Fatalln(v ...interface{}) {
	l.log(ERR, strings.TrimSuffix(fmt.Sprintln(v...), "\n"))
	os.Exit(1)
}

// Panic works like log.Panic. Use level ERR.
// Also output defaultKeyvals for prefixKeys/suffixKeys.
func (l *Logger) Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	l.log(ERR, s)
	panic(s)
}

// Panicf works like log.Panicf. Use level ERR.
// Also output defaultKeyvals for prefixKeys/suffixKeys.
func (l *Logger) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	l.log(ERR, s)
	panic(s)
}

// Panicln works like log.Panicln. Use level ERR.
// Also output defaultKeyvals for prefixKeys/suffixKeys.
func (l *Logger) Panicln(v ...interface{}) {
	s := strings.TrimSuffix(fmt.Sprintln(v...), "\n")
	l.log(ERR, s)
	panic(s)
}

func (l *Logger) log(level logLevel, msg interface{}, keyvals ...interface{}) {
	l.RLock()
	defer l.RUnlock()
	if l.parent != nil {
		l.RUnlock()
		l.mergeParent()
		l.RLock()
	}

	if *l.level > level {
		return
	}

	if len(keyvals)%2 != 0 {
		l.New().AddCallDepth(getPackageDepth()).Err("odd keyvals")
		keyvals = append(keyvals, MissingValue)
	}

	// TODO Combine all of this in single type and use sync.Pool.
	// Probably several different pools with different key sizes.
	// Use same len(keys) capability for all slices.
	// TODO Pre-calculate surroundKeys/prefixFormat/suffixFormat in
	// places where prefixKeys/suffixKeys may change.
	keys := make(map[string]interface{}, len(l.prefixKeys)+len(keyvals)/2+len(l.suffixKeys))
	prefixFormat := make([]string, 0, len(l.prefixKeys))
	suffixFormat := make([]string, 0, len(l.suffixKeys))
	middleFormat := make([]string, 0, len(keyvals)/2)
	middleKeys := make([]string, 0, len(keyvals)/2)
	values := make([]interface{}, 0, len(keys))
	surroundKeys := make(map[string]bool, len(l.prefixKeys)+len(l.suffixKeys))

	// Gather keys for output:
	// 1. Add prefixKeys which has non-nil defaultKeyVals.
	for _, k := range l.prefixKeys {
		surroundKeys[k] = true
		if l.defaultKeyvals[k] != nil {
			keys[k] = l.defaultKeyvals[k]
		}
		prefixFormat = append(prefixFormat, l.getFormat(k))
	}
	// 2. Add suffixKeys which has non-nil defaultKeyVals.
	for _, k := range l.suffixKeys {
		surroundKeys[k] = true
		if l.defaultKeyvals[k] != nil {
			keys[k] = l.defaultKeyvals[k]
		}
		suffixFormat = append(suffixFormat, l.getFormat(k))
	}
	// 3. Add msg to middleKeys. Msg value may be nil.
	middleKeys = append(middleKeys, KeyMessage)
	middleFormat = append(middleFormat, l.getFormat(KeyMessage))
	keys[KeyMessage] = msg
	// 4. Add keyvals to prefixKeys/middleKeys/suffixKeys.
	//    May overwrite prefixKeys/suffixKeys values from defaultKeyvals.
	//    May have nil values.
	for i := 0; i < len(keyvals); i += 2 {
		k, ok := keyvals[i].(string)
		if !ok {
			l.New().AddCallDepth(getPackageDepth()).SetKeyValFormat(" %#[2]v").Err("key is not string", "key", keyvals[i])
			k = fmt.Sprint(keyvals[i])
		}
		if !surroundKeys[k] {
			middleKeys = append(middleKeys, k)
			middleFormat = append(middleFormat, l.getFormat(k))
		}
		if t, ok := keyvals[i+1].(time.Time); ok {
			keys[k] = t.Format(*l.timeValFormat)
		} else {
			keys[k] = keyvals[i+1]
		}
	}
	// 5. Add current time if output format is JSON.
	if *l.format == JSON {
		keys[KeyTime] = time.Now().UTC().Format(*l.timeFormat)
	}
	// 6. Add log level.
	keys[KeyLevel] = level
	// 7. Add unit unless user set it to nil.
	//    If user didn't provide custom value then use package name.
	unit, okUnit := keys[KeyUnit]
	// 8. Add func and source unless user set them to nil.
	_, okFunc := keys[KeyFunc]
	_, okSource := keys[KeySource]
	if okUnit && unit == Auto || okSource || okFunc {
		if pc, file, line, ok := runtime.Caller(l.callDepth); ok {
			dir, file := path.Split(file)
			if okUnit && unit == Auto {
				keys[KeyUnit] = path.Base(dir)
			}
			if okFunc {
				keys[KeyFunc] = path.Base(runtime.FuncForPC(pc).Name())
			}
			if okSource {
				keys[KeySource] = fmt.Sprintf("%s:%d", file, line)
			}
		}
	}
	// 9. Add stack trace if user asks for it.
	//    If user didn't provide custom value then use default one.
	stack, okStack := keys[KeyStack]
	if okStack && stack == Auto {
		const size = 64 << 10
		buf := make([]byte, size)
		keys[KeyStack] = string(buf[:runtime.Stack(buf, false)])
	}

	// Now we've prepared all middleKeys plus some prefixKeys/suffixKeys
	// which wasn't disabled (nil in defaultKeyvals) and was provided
	// (non-nil in defaultKeyvals or anything in keyvals) by user.
	for i, k := range l.prefixKeys {
		if _, ok := keys[k]; ok {
			values = append(values, fmt.Sprintf(prefixFormat[i], k, keys[k]))
		}
	}
	for i, k := range middleKeys {
		values = append(values, fmt.Sprintf(middleFormat[i], k, keys[k]))
	}
	for i, k := range l.suffixKeys {
		if _, ok := keys[k]; ok {
			values = append(values, fmt.Sprintf(suffixFormat[i], k, keys[k]))
		}
	}

	// Output.
	if *l.format == Text {
		log.Print(values...)
	} else {
		// TODO Split this function into separate ones for Text
		// and JSON formats, to avoid useless text formatting for JSON.
		buf, err := json.Marshal(keys)
		if err != nil {
			log.Print(err)
		} else {
			log.Print(string(buf))
		}
	}
}

// mergeParent will merge l.parent's settings into l.
//
// mergeParent should be called in lazy way before using l settings.
// This allows to apply default configuration changes done in main() on
// DefaultLogger to package-global log vars in other packages, if they
// didn't already used log within init().
//
//   format:         use parent only by default
//   level:          use parent only by default
//   keyValFormat:   use parent only by default
//   timeFormat:     use parent only by default
//   timeValFormat:  use parent only by default
//   callDepth:      add parent's
//   defaultKeyvals: use parent only by default (set key to nil to drop parent's value)
//   prefixKeys:     prepend parent's keys (XXX no ease way to replace!)
//   suffixKeys:     append  parent's keys (XXX no ease way to replace!)
//   keysFormat:     use parent only by default (set to DefaultKeyValFormat to drop parent's value)
func (l *Logger) mergeParent() {
	// Handle recursive calls, like in case "key is not string".
	l.RLock()
	if l.parent == nil {
		l.RUnlock()
		return
	}
	l.RUnlock()

	l.Lock()
	defer l.Unlock()
	p := l.parent
	if p == nil {
		return
	}
	p.mergeParent()
	p.RLock()
	defer p.RUnlock()

	if l.format == nil {
		l.format = p.format
	}
	if l.level == nil {
		l.level = p.level
	}
	if l.keyValFormat == nil {
		l.keyValFormat = p.keyValFormat
	}
	if l.timeFormat == nil {
		l.timeFormat = p.timeFormat
	}
	if l.timeValFormat == nil {
		l.timeValFormat = p.timeValFormat
	}
	l.callDepth += p.callDepth
	for k, v := range p.defaultKeyvals {
		if _, ok := l.defaultKeyvals[k]; !ok {
			l.defaultKeyvals[k] = v
		}
	}
	l.prefixKeys = append(append([]string(nil), p.prefixKeys...), l.prefixKeys...)
	l.suffixKeys = append(l.suffixKeys, p.suffixKeys...)
	for k, v := range p.keysFormat {
		if _, ok := l.keysFormat[k]; !ok {
			l.keysFormat[k] = v
		}
	}

	l.parent = nil
}

// getFormat returns keyValFormat for k.
//
// mergeParent must be called before getFormat.
func (l *Logger) getFormat(k string) string {
	if format, ok := l.keysFormat[k]; ok {
		return format
	}
	return *l.keyValFormat
}

// getPackageDepth returns current stack depth within caller's package.
func getPackageDepth() int {
	_, callerFile, _, ok := runtime.Caller(1)
	callerPkg, _ := path.Split(callerFile)
	for depth := 0; ok; depth++ {
		var nextFile string
		_, nextFile, _, ok = runtime.Caller(2 + depth)
		nextPkg, _ := path.Split(nextFile)
		if callerPkg != nextPkg {
			return depth
		}
	}
	return 0
}

// getErr returns first arg of type error or msg.
func getErr(msg interface{}, keyvals ...interface{}) error {
	if err, ok := msg.(error); ok {
		return err
	}
	for _, keyval := range keyvals {
		if err, ok := keyval.(error); ok {
			return err
		}
	}
	return fmt.Errorf("%s", msg)
}
