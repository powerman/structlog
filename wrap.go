package structlog

type keyvalsError struct {
	err     error
	keyvals []interface{}
}

// Error implements error interface.
func (err *keyvalsError) Error() string {
	return err.err.Error()
}

// Cause implements github.com/pkg/errors.causer interface.
func (err *keyvalsError) Cause() error {
	return err.err
}

// Unwrap implements interface used by errors.Unwrap.
func (err *keyvalsError) Unwrap() error {
	return err.err
}

func unwrap(err error) (keyvals []interface{}) {
	for err != nil {
		switch errWith := err.(type) { //nolint:errorlint // Needs to also support Cause.
		case *keyvalsError:
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

	return &keyvalsError{
		err:     err,
		keyvals: keyvals,
	}
}
