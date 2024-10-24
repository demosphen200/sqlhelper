package ex

import "fmt"

type Log struct {
	indentStepLength   int
	currentIndentLevel int
	currentIndent      string
}

func (log *Log) createIndent() string {
	if log.indentStepLength <= 0 {
		log.indentStepLength = 4
	}
	indentLength := log.indentStepLength * log.currentIndentLevel
	indent := make([]byte, indentLength)
	for t := 0; t < indentLength; t++ {
		indent[t] = ' '
	}
	return string(indent)
}

func (log *Log) IncIndent() *Log {
	log.currentIndentLevel++
	log.currentIndent = log.createIndent()
	return log
}

func (log *Log) DecIndent() *Log {
	if log.currentIndentLevel > 0 {
		log.currentIndentLevel--
		log.currentIndent = log.createIndent()
	}
	return log
}

func (log *Log) Printf(format string, args ...any) *Log {
	fmt.Printf(
		fmt.Sprintf("%s%s", log.currentIndent, format),
		args...,
	)
	return log
}

func (log *Log) Printlnf(format string, args ...any) *Log {
	fmt.Printf(
		fmt.Sprintf("%s%s\n", log.currentIndent, format),
		args...,
	)
	return log
}

func (log *Log) WithIndent(fn func()) *Log {
	log.IncIndent()
	fn()
	log.DecIndent()
	return log
}

func (log *Log) Println(line string) *Log {
	fmt.Printf("%s%s\n", log.currentIndent, line)
	return log
}
