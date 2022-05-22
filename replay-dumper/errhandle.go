package main

import (
	"log"
	"path"
	"runtime"
)

//lint:ignore U1000 must
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

//lint:ignore U1000 must
func fmust(err error) {
	if err != nil {
		pc, filename, line, _ := runtime.Caller(1)
		log.Fatalf("Error: %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), path.Base(filename), line, err)
	}
}

func noerr[T any](t T, err error) T {
	if err != nil {
		pc, filename, line, _ := runtime.Caller(1)
		log.Fatalf("Error: %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), path.Base(filename), line, err)
	}
	return t
}

func noerr2[T any, T1 any](t T, t1 T1, err error) (T, T1) {
	if err != nil {
		pc, filename, line, _ := runtime.Caller(1)
		log.Fatalf("Error: %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), path.Base(filename), line, err)
	}
	return t, t1
}
