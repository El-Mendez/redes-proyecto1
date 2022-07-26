package utils

// Successful is a utility function that logs and stops the program when it receives something that is not nil. It is
// expected to use this with a function that returns an error or nil.
func Successful(err error, format string) {
	if err != nil {
		Logger.Fatalf(format, err)
	}
}
