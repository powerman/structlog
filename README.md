# structlog
[![GoDoc](https://godoc.org/github.com/powerman/structlog?status.svg)](http://godoc.org/github.com/powerman/structlog) [![Go Report Card](https://goreportcard.com/badge/github.com/powerman/structlog)](https://goreportcard.com/report/github.com/powerman/structlog) [![CircleCI](https://circleci.com/gh/powerman/structlog.svg?style=svg)](https://circleci.com/gh/powerman/structlog) [![Coverage Status](https://coveralls.io/repos/powerman/structlog/badge.svg?branch=master&service=github)](https://coveralls.io/github/powerman/structlog?branch=master)

## Features

- log only key/value pairs
- output both as Text and JSON
- log level support
- compatible enough with log.Logger to use as drop-in replacement
- short names for service keys (like log level, time, etc.)
- support default values for keys
- service keys with caller's function name, file and line
- fixed log level width
- service key for high-level source name (like package or subsystem),
  caller's package name by default
- when creating new logger instance:
    - inherit default logger settings (configured in main())
    - add new default key/values
    - disable inherited default keys
- warn about imbalanced key/value pairs
- first parameter to log functions should be value for "message" service key
- able to output stack trace
- level-guards like IsDebug()
- output complex struct as key values (using "%v" like formatting)
- Error returns message as error (auto-convert from string, if needed)
    - actually it returns first `.(error)` arg if any or message otherwise
- convenient helpers IfFail and Recover for use with defer
- output can be redirected/intercepted
- when output as JSON:
    - add service field time by default
- when output as Text:
    - do not add service field time by default
    - order of keys in output is fixed, same as order of log function
      parameters
    - it is possible to set minimal width for some keys (both user and
      service ones)
    - it is possible to setup order for service keys, including this keys
      should be output before and after user keys
    - it is possible to choose output style for some keys: "key: ",
      "key=", do not output key name
    - it is possible to choose how to escape values by default and for
      selected keys: no escaping, use \`\`, use ""
    - do not support colored output - this can be added by external tools
      like grc and we anyway won't have colors in log files
