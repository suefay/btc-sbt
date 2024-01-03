package statemachine

// ErrKind defines the error kind
type ErrKind uint8

const (
	InvalidOpErr ErrKind = iota
	ExecutionFailedErr
)

// InvalidOperation defines the error due to invalid operation when transitioning the state
type InvalidOperation struct {
	reason string
}

// NewInvalidOperation creates an InvalidOperation instance
func NewInvalidOperation(reason string) InvalidOperation {
	return InvalidOperation{reason}
}

// Error implements the error interface
func (io InvalidOperation) Error() string {
	return io.reason
}

// ExecutionFailed defines the error due to execution failed when transitioning the state
type ExecutionFailed struct {
	reason string
}

// NewExecutionFailed creates an ExecutionFailed instance
func NewExecutionFailed(reason string) ExecutionFailed {
	return ExecutionFailed{reason}
}

// Error implements the error interface
func (ef ExecutionFailed) Error() string {
	return ef.reason
}

// wrapError wraps an error to the corresponding error type
func wrapError(kind ErrKind, err error) error {
	if err != nil {
		switch kind {
		case InvalidOpErr:
			return NewInvalidOperation(err.Error())

		case ExecutionFailedErr:
			return NewExecutionFailed(err.Error())

		default:
			return err
		}
	}

	return nil
}

// IsInvalidOpErr checks if the given error is InvalidOperation
func IsInvalidOpErr(err error) bool {
	_, ok := err.(InvalidOperation)
	return ok
}

// IsExecutionFailedErr checks if the given error is ExecutionFailed
func IsExecutionFailedErr(err error) bool {
	_, ok := err.(ExecutionFailed)
	return ok
}
