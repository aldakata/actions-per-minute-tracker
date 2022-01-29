package main

import (
	"fmt"
	"log"
	"syscall"
	"time"

	"actions-per-minute-tracker/win32"
)

type callback struct {
	hookKeyboard win32.HHOOK
	hookMouse    win32.HHOOK
}

func (c *callback) keyboardCallback(code int, wparam win32.WPARAM, lparam win32.LPARAM) win32.LRESULT {
	if code >= 0 {
		if wparam == win32.WM_KEYDOWN {
			fmt.Println("got called keyboard", wparam)
		}
	}
	return win32.CallNextHookEx(c.hookKeyboard, code, wparam, lparam)
}

func (c *callback) mouseCallback(code int, wparam win32.WPARAM, lparam win32.LPARAM) win32.LRESULT {
	if code >= 0 {
		if wparam == win32.WM_LBUTTONDOWN ||
			wparam == win32.WM_RBUTTONDOWN ||
			wparam == win32.WM_XBUTTONDOWN ||
			wparam == win32.WM_MBUTTONDOWN {
			fmt.Println("got called mouse", wparam)
		} else {
		}
	}
	return win32.CallNextHookEx(c.hookMouse, code, wparam, lparam)
}

func windowProc(hwnd syscall.Handle, msg uint32, wparam win32.WPARAM, lparam win32.LPARAM) win32.LRESULT {
	return win32.DefWindowProc(hwnd, msg, wparam, lparam)
}

func main() {
	fmt.Println("starting main...")
	// setup apm tracker
	cb := callback{}
	hookKeyboard := win32.SetWindowsHookEx(win32.WH_KEYBOARD_LL, cb.keyboardCallback, 0, 0)
	hookMouse := win32.SetWindowsHookEx(win32.WH_MOUSE_LL, cb.mouseCallback, 0, 0)

	defer func() {
		win32.UnhookWindowsHookEx(hookKeyboard)
		win32.UnhookWindowsHookEx(hookMouse)
	}()

	cb.hookKeyboard = hookKeyboard
	cb.hookMouse = hookMouse

	// setup window
	className := "apm-window-object"

	instance, err := win32.GetModuleHandle()
	if err != nil {
		log.Println(err)
		return
	}

	cursor, err := win32.LoadCursorResource(win32.IDC_ARROW)
	if err != nil {
		log.Println(err)
		return
	}

	wndClass := win32.NewWNDClasss(className, windowProc, instance, cursor)
	if _, err = win32.RegisterClassEx(&wndClass); err != nil {
		log.Println(err)
		return
	}

	height := 35
	width := 130
	var extraStyles uint32 = win32.WS_EX_COMPOSITED | win32.WS_EX_LAYERED | win32.WS_EX_NOACTIVATE | win32.WS_EX_TOPMOST | win32.WS_EX_TRANSPARENT
	var styles uint32 = win32.WS_VISIBLE | win32.WS_POPUP
	hwnd := win32.CreateWindow(
		extraStyles, // extra style
		className,
		"Test Pupper Window", // name
		uint32(styles),       // style
		int64(win32.GetSystemMetrics(win32.SM_CXSCREEN))-int64(width), // x
		int64(height*2), // y
		int64(width),    // width
		int64(height),   // height
		0,               // parent
		0,               // menu
		instance,
	)

	win32.SetWindowPos(
		hwnd,
		^win32.HWND(0),
		0,
		0,
		0,
		0,
		win32.SWP_NOACTIVATE|win32.SWP_NOMOVE|win32.SWP_NOSIZE|win32.SWP_SHOWWINDOW,
	)

	counter := 1000
	ticker := time.NewTicker(1000 * time.Millisecond)
	go func() {
		for {
			select {
			case <-ticker.C:
				win32.SendMessage(hwnd, 35001, 0, 0)
				counter++
			}
		}
	}()

	// pull messages
	var msg win32.MSG
	for {
		fmt.Println("looping...")
		msgVal := win32.GetMessage(&msg, 0, 0, 0)
		if msgVal <= 0 {
			fmt.Println("bad msg val", msgVal)
			break
		}
		fmt.Println("got msg", msg)
		win32.TranslateMessage(&msg)
		win32.DispatchMessage(&msg)
	}
}