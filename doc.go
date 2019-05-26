// Package structlog provides structured logger which looks friendly as
// plain text (more like handcrafted vertical-aligned log lines which
// doesn't feels like key/value pairs) but also really useful as JSON
// (each important value in log line became separate key/value pair).
//
// In short, with structlog you will have to write just this:
//
//	log := structlog.New()
//	log.Warn("something goes wrong", "somevar", var, "err", err)
//
// and structlog will take care of the rest (output plain text or JSON,
// format plain text output to make it easy to read, include in output
// extra key/value pairs - both added by you before that call and
// calculated automatically like caller's details or stack trace).
//
// Well, not really. First you'll need to configure logger, to define what
// "easy to read" means in your case - but this has to be done just once,
// usually in main(), and then you'll get simple and powerful logging as
// promised above.
//
// Overview
//
// It was designed to produce easy to read, vertically-aligned log lines
// while your project is small, and later, when your project will grow to
// use something like ELK Stack, switch by changing one option to
// produce easy to parse JSON log records.
//
// You can log key/value pairs without bothering about output format
// (plain text or JSON) and then fine-tune it plain text output by:
//
//   - Re-ordering output values by choosing which keys should be output
//     as prefix (before log message) and as suffix (after other keys).
//   - Defining internal order of prefix and suffix keys.
//   - Defining key/value output style using fmt.Sprintf 'verbs', including
//     ability to output only value, without key name.
//   - Choosing which ones of pre-defined key/values to include in output:
//     - caller's package
//     - caller's function name
//     - caller's file and line
//     - multiline stack trace (a-la panic output)
//
// Supported log levels: Err, Warn, Info and Debug.
//
// On import it calls stdlib's log.SetFlags(0) and by default will use
// stdlib's log.Print() to output log lines - this is to make sure
// structlog's output goes at same place as logging from other packages
// (which often use stdlib's log).
//
// Inheritance
//
// For convenience you may have multiple loggers, and create new logger in
// such a way to inherit all settings from existing logger. These settings
// often include predefined key/value pairs to be automatically output in
// each logged line in addition to logger formatting setting.
//
// All loggers created using structlog.New() will inherit settings from
// logger in global variable structlog.DefaultLogger. So, usually you will
// configure structlog.DefaultLogger in your main(), to apply this
// configuration to each and every logger created in any other package of
// your application.
//
// Then you can do some extra setup on your logger (usually - just add
// some default key/value pair which should be output by each log call),
// and call log.New() method to get new logger which will inherit this
// extra setup.
//
// E.g. imagine HTTP middlewares: first will detect IP of connected client
// and store it in logger, second will check authentication and optionally
// add key/value pair with user ID - as result everything logged by your
// HTTP handler later using logger preconfigured by these middlewares will
// include remote IP and user ID in each log record - without needs to
// manually include it in each line where you log something.
//
// Contents
//
// ★ Creating a new logger:
//
//	New (function)  - will inherit from structlog.DefaultLogger
//	New (method)    - will inherit from it's object
//	NewZeroLogger   - new empty logger (usually you won't need this)
//
// ★ Passing logger inside context.Context:
//
//	NewContext
//	FromContext
//
// ★ Normal logging:
//
//	Debug
//	Info
//	Warn
//	Err
//	PrintErr        - like Err, but don't return error (usually you won't need this)
//
// ★ Logging useful with defer:
//
//	DebugIfFail
//	InfoIfFail
//	WarnIfFail
//	ErrIfFail
//	Recover
//
// ★ Delayed logging:
//
//	WrapErr
//
// ★ Configuring structlog.DefaultLogger in your main():
//
//	AppendPrefixKeys
//	PrependSuffixKeys
//	SetKeyValFormat
//	SetKeysFormat
//	SetLogFormat
//	SetPrefixKeys
//	SetSuffixKeys
//	SetTimeFormat
//	SetTimeValFormat
//
// ★ Configuring current logger:
//
//	SetDefaultKeyvals
//	AddCallDepth
//
// ★ Handling log levels:
//
//	IsDebug
//	IsInfo
//	ParseLevel
//	SetLogLevel
//
// ★ Passing this logger to 3rd-party packages which expects interface of stdlib's log.Logger:
//
//	Fatal
//	Fatalf
//	Fatalln
//	Panic
//	Panicf
//	Panicln
//	Print
//	Printf
//	Println
//
// ★ Redirecting log output (useful to redirect to ioutil.Discard in tests):
//
//	SetOutput
//	SetPrinter
package structlog
