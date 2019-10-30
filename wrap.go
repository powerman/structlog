package structlog

type errWithKeyvals struct {
	err     error
	keyvals []interface{}
}

// Error implements error interface.
func (err *errWithKeyvals) Error() string {
	return err.err.Error()
}

// Cause implements github.com/pkg/errors.causer interface.
func (err *errWithKeyvals) Cause() error {
	return err.err
}

// Unwrap implements interface used by errors.Unwrap.
func (err *errWithKeyvals) Unwrap() error {
	return err.err
}

func unwrap(err error) (keyvals []interface{}) {
	for err != nil {
		switch errWith := err.(type) {
		case *errWithKeyvals:
			keyvals = append(errWith.keyvals, keyvals...)
			err = errWith.Unwrap()
		case interface{ Unwrap() error }:
			err = errWith.Unwrap()
		case interface{ Cause() error }:
			err = errWith.Cause()
		default:
			err = nil
		}
	}
	return keyvals
}

// WrapErr returns given err wrapped with keyvals. If returned err will be
// logged later these keyvals will be included in output.
//
// If called with nil error it'll return nil.
func (l *Logger) WrapErr(err error, keyvals ...interface{}) error {
	if err == nil {
		return nil
	}

	if len(keyvals)%2 != 0 {
		l.New().AddCallDepth(getPackageDepth()).PrintErr("odd keyvals")
		keyvals = append(keyvals, MissingValue)
	}

	return &errWithKeyvals{
		err:     err,
		keyvals: keyvals,
	}
}
