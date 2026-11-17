#!/usr/bin/env python3
"""Park or report the virtual display's real X pointer.

Usage: x-pointer.py park | x-pointer.py where

Xvfb starts its pointer at the centre of the screen, where the browser
windows are. When another test worker opens or closes a window under it,
the X server sends that position to the window below as a mouse move, in
the middle of whatever gesture a test is making there. `park` moves the
pointer to the screen's bottom-right corner, outside every browser window.
`where` prints the pointer's position as "x,y".

Uses libX11 through ctypes, so the image needs no extra package.
"""
import ctypes
import sys

X = ctypes.cdll.LoadLibrary("libX11.so.6")
X.XOpenDisplay.argtypes = [ctypes.c_char_p]
X.XOpenDisplay.restype = ctypes.c_void_p
X.XDefaultScreen.argtypes = [ctypes.c_void_p]
X.XDisplayWidth.argtypes = [ctypes.c_void_p, ctypes.c_int]
X.XDisplayHeight.argtypes = [ctypes.c_void_p, ctypes.c_int]
X.XDefaultRootWindow.argtypes = [ctypes.c_void_p]
X.XDefaultRootWindow.restype = ctypes.c_ulong
X.XWarpPointer.argtypes = [ctypes.c_void_p, ctypes.c_ulong, ctypes.c_ulong, ctypes.c_int, ctypes.c_int,
                           ctypes.c_uint, ctypes.c_uint, ctypes.c_int, ctypes.c_int]
X.XQueryPointer.argtypes = [ctypes.c_void_p, ctypes.c_ulong] + [ctypes.c_void_p] * 7
X.XSync.argtypes = [ctypes.c_void_p, ctypes.c_int]
X.XCloseDisplay.argtypes = [ctypes.c_void_p]


def corner(display):
    screen = X.XDefaultScreen(display)
    return X.XDisplayWidth(display, screen) - 1, X.XDisplayHeight(display, screen) - 1


def where(display, root):
    window, child = ctypes.c_ulong(), ctypes.c_ulong()
    root_x, root_y, win_x, win_y = (ctypes.c_int() for _ in range(4))
    mask = ctypes.c_uint()
    X.XQueryPointer(display, root, ctypes.byref(window), ctypes.byref(child), ctypes.byref(root_x),
                    ctypes.byref(root_y), ctypes.byref(win_x), ctypes.byref(win_y), ctypes.byref(mask))
    return root_x.value, root_y.value


def main(command):
    display = X.XOpenDisplay(None)
    if not display:
        sys.exit("x-pointer: cannot open the display")
    root = X.XDefaultRootWindow(display)
    if command == "park":
        x, y = corner(display)
        X.XWarpPointer(display, 0, root, 0, 0, 0, 0, x, y)
        X.XSync(display, 0)
    print("%d,%d" % where(display, root))
    X.XCloseDisplay(display)


if __name__ == "__main__":
    if len(sys.argv) != 2 or sys.argv[1] not in ("park", "where"):
        sys.exit("usage: x-pointer.py park | where")
    main(sys.argv[1])
