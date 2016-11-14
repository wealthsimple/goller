package goller

import "log"

func checkErr(e error, l *log.Logger) {
	if e != nil {
		l.Fatal(e)
	}
}
