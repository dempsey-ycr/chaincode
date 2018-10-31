package logging

import "github.com/hyperledger/fabric/core/chaincode/shim"

var log *shim.ChaincodeLogger

// Debug logs will only appear if the ChaincodeLogger LoggingLevel is set to
// LogDebug.
func Debug(args ...interface{}) {
	log.Debug(args)
}

// Info logs will appear if the ChaincodeLogger LoggingLevel is set to
// LogInfo or LogDebug.
func Info(args ...interface{}) {
	log.Info(args)
}

// Notice logs will appear if the ChaincodeLogger LoggingLevel is set to
// LogNotice, LogInfo or LogDebug.
func Notice(args ...interface{}) {
	log.Notice(args)
}

// Warning logs will appear if the ChaincodeLogger LoggingLevel is set to
// LogWarning, LogNotice, LogInfo or LogDebug.
func Warning(args ...interface{}) {
	log.Warning(args)
}

// Error logs will appear if the ChaincodeLogger LoggingLevel is set to
// LogError, LogWarning, LogNotice, LogInfo or LogDebug.
func Error(args ...interface{}) {
	log.Error(args)
}

// Critical logs always appear; They can not be disabled.
func Critical(args ...interface{}) {
	log.Critical(args)
}

// Debugf logs will only appear if the ChaincodeLogger LoggingLevel is set to
// LogDebug.
func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

// Infof logs will appear if the ChaincodeLogger LoggingLevel is set to
// LogInfo or LogDebug.
func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

// Noticef logs will appear if the ChaincodeLogger LoggingLevel is set to
// LogNotice, LogInfo or LogDebug.
func Noticef(format string, args ...interface{}) {
	log.Noticef(format, args...)
}

// Warningf logs will appear if the ChaincodeLogger LoggingLevel is set to
// LogWarning, LogNotice, LogInfo or LogDebug.
func Warningf(format string, args ...interface{}) {
	log.Warningf(format, args...)
}

// Errorf logs will appear if the ChaincodeLogger LoggingLevel is set to
// LogError, LogWarning, LogNotice, LogInfo or LogDebug.
func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

// Criticalf logs always appear; They can not be disabled.
func Criticalf(format string, args ...interface{}) {
	log.Criticalf(format, args...)
}

func init() {
	if log == nil {
		log = shim.NewLogger("verse")
		log.SetLevel(shim.LogInfo)
	}
}
