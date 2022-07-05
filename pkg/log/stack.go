package log

import (
	"fmt"
	"path"
	"runtime"
	"strconv"
	"strings"
)

// Stack represents a Stack of program counters.
type Stack []uintptr

type Frame uintptr

func (f Frame) pc() uintptr { return uintptr(f) - 1 }

func Callers(skip int) *Stack {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(skip, pcs[:])
	var st Stack = pcs[0:n]
	return &st
}

// file returns the full path to the file that contains the
// function for this Frame's pc.
func (f Frame) file() string {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return "unknown"
	}
	file, _ := fn.FileLine(f.pc())
	return file
}

// line returns the line number of source code of the
// function for this Frame's pc.
func (f Frame) line() int {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return 0
	}
	_, line := fn.FileLine(f.pc())
	return line
}

func (f Frame) name() string {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return "unknown"
	}
	return fn.Name()
}

func (f Frame) String() string {
	dir, file := path.Split(f.file())
	pkg := path.Base(dir)
	return path.Join(pkg, file) + ":" + strconv.Itoa(f.line())
	//return path.TaskBase(f.file()) + ":" + strconv.Itoa(f.line())
}

func (s *Stack) Format(st fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case st.Flag('+'):
			for _, pc := range *s {
				f := Frame(pc)
				fmt.Fprintf(st, "\n%+v", f)
			}
		}
	}
}

// Top 需要在日志打印函数中间接调用确保获取的调用栈层数正确
func Top() Frame {
	st := Callers(3)
	for _, f := range *st {
		return Frame(f)
	}
	return 0
}

func funcname(name string) string {
	i := strings.LastIndex(name, "/")
	name = name[i+1:]
	i = strings.Index(name, ".")
	return name[i+1:]
}
