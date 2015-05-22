// +build linux,!gccgo

package main

/*
#cgo CFLAGS: -Wall
extern void joinmnt();
void __attribute__((constructor)) initmnt(void) {
	joinmnt();
}
*/
import "C"
