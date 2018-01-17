package main

/*
#cgo LDFLAGS: -ldl

#include <dlfcn.h>
#include <stdlib.h>

void call_SACLockScreenImmediate(void *addr) {
	void (*fn)(void) = addr;
	fn();
}
*/
import "C"

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"unsafe"

	"github.com/shurcooL/trayhost"
)

const (
	framework = "/System/Library/PrivateFrameworks/login.framework/Versions/Current/login"
	symbol    = "SACLockScreenImmediate"
)

var (
	funcAddr unsafe.Pointer
)

func init() {
	s := C.CString(framework)
	defer C.free(unsafe.Pointer(s))

	handle := C.dlopen(s, C.int(C.RTLD_NOW))
	if handle == nil {
		log.Fatal("error opening framework")
	}

	sym := C.CString(symbol)
	defer C.free(unsafe.Pointer(sym))
	symhandle := C.dlsym(handle, sym)
	if symhandle == nil {
		log.Fatal("error resolving symbol")
	}

	funcAddr = unsafe.Pointer(symhandle)
}

func main() {
	log.Println("started")

	menuItems := []trayhost.MenuItem{
		{
			Title: "Lock Screen",
			Handler: func() {
				C.call_SACLockScreenImmediate(funcAddr)
			},
		},
		trayhost.SeparatorMenuItem(),
		{
			Title:   "Quit",
			Handler: trayhost.Exit,
		},
	}

	// On macOS, when you run an app bundle, the working directory of the executed process
	// is the root directory (/), not the app bundle's Contents/Resources directory.
	// Change directory to Resources so that we can load resources from there.
	ep, err := os.Executable()
	if err != nil {
		log.Fatal("os.Executable:", err)
	}
	err = os.Chdir(filepath.Join(filepath.Dir(ep), "..", "Resources"))
	if err != nil {
		log.Fatal("os.Chdir:", err)
	}

	// Load tray icon.
	iconData, err := ioutil.ReadFile("lock.png")
	if err != nil {
		log.Fatal(err)
	}

	trayhost.Initialize("traylock", iconData, menuItems)
	trayhost.EnterLoop()
}
