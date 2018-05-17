package log

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time"
)

const (
	infoType  = "\033[32mINFO \033[0m"
	warnType  = "\033[33mWARN \033[0m"
	errType   = "\033[31mERR  \033[0m"
	panicType = "\033[37;41mPANIC\033[0m"
	traceType = "\033[36mTRACE\033[0m"
)

func blue(val interface{}) string {
	return fmt.Sprintf("\033[34m%v\033[0m", val)
}

func log(typ string, msg string, vals ...interface{}) {
	var valsStr string
	for i := 0; i < (len(vals) - 1); i += 2 {
		valsStr += fmt.Sprintf("%v=%v ", blue(vals[i]), vals[i+1])
	}

	_, file, line, _ := runtime.Caller(2)
	a, _ := filepath.Abs(".")
	f, _ := filepath.Rel(a, file)
	f = filepath.FromSlash("./") + f
	callerStr := fmt.Sprintf("%v:%v", f, line)

	fmt.Printf("%v [%v] %v\t\t%v\t\t%v\n", typ, time.Now().Format("15:04:05.000"), msg, valsStr, callerStr)
}

// Info logs an information message.
func Info(msg string, vals ...interface{}) {
	log(infoType, msg, vals...)
}

// Warn logs a warning message.
func Warn(msg string, vals ...interface{}) {
	log(warnType, msg, vals...)
}

// Error logs an error message.
func Error(msg string, vals ...interface{}) {
	log(errType, msg, vals...)
}

// Panic logs a panic message and panics.
func Panic(msg string, vals ...interface{}) {
	log(panicType, msg, vals...)
	panic(msg)
}

// Trace prints the calling function name.
func Trace() {
	fn, _, _, _ := runtime.Caller(1)
	name := runtime.FuncForPC(fn).Name()
	name = filepath.Base(name)
	log(traceType, name+"()")
}
