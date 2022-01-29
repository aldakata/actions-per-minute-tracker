package win32

import (
	"syscall"
	"unsafe"
)

var (
	user32                  = syscall.NewLazyDLL("user32.dll")
	procSetWindowsHookEx    = user32.NewProc("SetWindowsHookExW")
	procUnhookWindowsHookEx = user32.NewProc("UnhookWindowsHookEx")
	procGetMessage          = user32.NewProc("GetMessageW")
	procCallNextHookEx      = user32.NewProc("CallNextHookEx")
	pLoadCursorW            = user32.NewProc("LoadCursorW")

	kernel32             = syscall.NewLazyDLL("kernel32.dll")
	procGetModuleHandleW = kernel32.NewProc("GetModuleHandleW")

	pDefWindowProcW   = user32.NewProc("DefWindowProcW")
	pRegisterClassExW = user32.NewProc("RegisterClassExW")
	pCreateWindowExW  = user32.NewProc("CreateWindowExW")
	pGetSystemMetrics = user32.NewProc("GetSystemMetrics")
	pSetWindowPos     = user32.NewProc("SetWindowPos")
	pSendMessage      = user32.NewProc("SendMessageW")

	pTranslateMessage = user32.NewProc("TranslateMessage")
	pDispatchMessageW = user32.NewProc("DispatchMessageW")
)

const (
	IDC_ARROW    = 32512
	COLOR_WINDOW = 5
)

// https://docs.microsoft.com/en-us/windows/win32/inputdev/keyboard-input-notifications
const (
	WM_KEYDOWN = 256
)

// https://docs.microsoft.com/en-us/windows/win32/inputdev/mouse-input-notifications
const (
	WM_LBUTTONDOWN = 513
	WM_RBUTTONDOWN = 516
	WM_MBUTTONDOWN = 519
	WM_XBUTTONDOWN = 523
)

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-setwindowshookexw
const (
	WH_KEYBOARD_LL = 13
	WH_MOUSE_LL    = 14
)

// https://docs.microsoft.com/en-us/windows/win32/winprog/windows-data-types
type (
	WPARAM  uintptr
	LPARAM  uintptr
	LRESULT uintptr
	HHOOK   uintptr
)

// https://docs.microsoft.com/en-us/windows/win32/winmsg/lowlevelkeyboardproc
type HOOKPROC func(int, WPARAM, LPARAM) LRESULT

// http://msdn.microsoft.com/en-us/library/windows/desktop/ms644958.aspx
type MSG struct {
	Hwnd    uintptr
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      POINT
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/dd162805.aspx
type POINT struct {
	X, Y int32
}

// GetSystemMetrics constants
const (
	SM_CXSCREEN = 0
	SM_CYSCREEN = 1
)

// WindowFlags
const (
	// http://msdn.microsoft.com/en-us/library/windows/desktop/ff700543(v=vs.85).aspx
	WS_EX_TOPMOST     = 8
	WS_EX_TRANSPARENT = 32
	WS_EX_LAYERED     = 0x00080000
	WS_EX_COMPOSITED  = 0x02000000
	WS_EX_NOACTIVATE  = 0x08000000

	// http://msdn.microsoft.com/en-us/library/windows/desktop/ms632600(v=vs.85).aspx
	WS_DISABLED = 134217728
	WS_VISIBLE  = 268435456
	WS_POPUP    = 0x80000000
)

const (
	SWP_NOSIZE       = 0x0001
	SWP_NOMOVE       = 0x0002
	SWP_NOZORDER     = 0x0004
	SWP_NOREDRAW     = 0x0008
	SWP_NOACTIVATE   = 0x0010
	SWP_FRAMECHANGED = 0x0020
	SWP_SHOWWINDOW   = 0x0040
)

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/ns-winuser-wndclassexw
type WNDCLASSEXW struct {
	size       uint32
	style      uint32
	wndProc    uintptr
	clsExtra   int32
	wndExtra   int32
	instance   syscall.Handle
	icon       syscall.Handle
	cursor     syscall.Handle
	background syscall.Handle
	menuName   *uint16
	className  *uint16
	iconSm     syscall.Handle
}

type HWND uintptr

// https://docs.microsoft.com/en-us/windows/win32/learnwin32/writing-the-window-procedure
type WindowProcFn func(syscall.Handle, uint32, WPARAM, LPARAM) LRESULT

func NewWNDClasss(className string, windowProc WindowProcFn, instance, cursor syscall.Handle) WNDCLASSEXW {
	wnclass := WNDCLASSEXW{
		wndProc:    syscall.NewCallback(windowProc),
		instance:   instance,
		cursor:     cursor,
		background: COLOR_WINDOW + 1,
		className:  syscall.StringToUTF16Ptr(className),
	}
	wnclass.size = uint32(unsafe.Sizeof(wnclass))
	return wnclass
}

func DefWindowProc(hwnd syscall.Handle, msg uint32, wparam WPARAM, lparam LPARAM) LRESULT {
	ret, _, _ := pDefWindowProcW.Call(
		uintptr(hwnd),
		uintptr(msg),
		uintptr(wparam),
		uintptr(lparam),
	)
	return LRESULT(ret)
}

func RegisterClassEx(wcx *WNDCLASSEXW) (uint16, error) {
	ret, _, err := pRegisterClassExW.Call(
		uintptr(unsafe.Pointer(wcx)),
	)
	if ret == 0 {
		return 0, err
	}
	return uint16(ret), nil
}

func CreateWindow(exStyle uint32, className, windowName string, style uint32, x, y, width, height int64, parent, menu, instance syscall.Handle) HWND {
	ret, _, _ := pCreateWindowExW.Call(
		uintptr(0),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(className))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(windowName))),
		uintptr(style),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height),
		uintptr(parent),
		uintptr(menu),
		uintptr(instance),
		uintptr(0),
	)
	return HWND(ret)
}

func GetSystemMetrics(index int) int {
	ret, _, _ := pGetSystemMetrics.Call(
		uintptr(index))

	return int(ret)
}

func SetWindowPos(hwnd, hWndInsertAfter HWND, x, y, cx, cy int, uFlags uint) bool {
	ret, _, _ := pSetWindowPos.Call(
		uintptr(hwnd),
		uintptr(hWndInsertAfter),
		uintptr(x),
		uintptr(y),
		uintptr(cx),
		uintptr(cy),
		uintptr(uFlags))

	return ret != 0
}

func SendMessage(hwnd HWND, msg int, wparam WPARAM, lparam LPARAM) LRESULT {
	ret, _, _ := pSendMessage.Call(
		uintptr(hwnd),
		uintptr(msg),
		uintptr(wparam),
		uintptr(lparam),
	)
	return LRESULT(ret)
}

func DispatchMessage(msg *MSG) {
	pDispatchMessageW.Call(uintptr(unsafe.Pointer(msg)))
}

func TranslateMessage(msg *MSG) {
	pTranslateMessage.Call(uintptr(unsafe.Pointer(msg)))
}

func GetModuleHandle() (syscall.Handle, error) {
	ret, _, err := procGetModuleHandleW.Call(uintptr(0))
	if ret == 0 {
		return 0, err
	}
	return syscall.Handle(ret), nil
}

func LoadCursorResource(cursorName uint32) (syscall.Handle, error) {
	ret, _, err := pLoadCursorW.Call(
		uintptr(0),
		uintptr(uint16(cursorName)),
	)
	if ret == 0 {
		return 0, err
	}
	return syscall.Handle(ret), nil
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-setwindowshookexw
func SetWindowsHookEx(idHook int, lpfn HOOKPROC, hMod uintptr, dwThreadId uintptr) HHOOK {
	ret, _, _ := procSetWindowsHookEx.Call(
		uintptr(idHook),
		uintptr(syscall.NewCallback(lpfn)),
		uintptr(hMod),
		uintptr(dwThreadId),
	)
	return HHOOK(ret)
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-unhookwindowshookex
func UnhookWindowsHookEx(hhk HHOOK) bool {
	ret, _, _ := procUnhookWindowsHookEx.Call(
		uintptr(hhk),
	)
	return ret != 0
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getmessage
func GetMessage(msg *MSG, hwnd uintptr, msgFilterMin, msgFilterMax uint32) int {
	ret, _, _ := procGetMessage.Call(
		uintptr(unsafe.Pointer(msg)),
		uintptr(hwnd),
		uintptr(msgFilterMin),
		uintptr(msgFilterMax),
	)
	return int(ret)
}

// https://docs.microsoft.com/en-us/windows/win32/winmsg/lowlevelkeyboardproc
func CallNextHookEx(hhk HHOOK, nCode int, wparam WPARAM, lparam LPARAM) LRESULT {
	ret, _, _ := procCallNextHookEx.Call(
		uintptr(hhk),
		uintptr(nCode),
		uintptr(wparam),
		uintptr(lparam),
	)
	return LRESULT(ret)
}