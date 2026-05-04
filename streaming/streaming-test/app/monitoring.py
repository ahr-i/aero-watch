import ctypes
import sys
from ctypes import wintypes

from app.models import MonitorInfo


if sys.platform != "win32":
    raise SystemExit("This app currently supports Windows only.")


user32 = ctypes.windll.user32


class RECT(ctypes.Structure):
    _fields_ = [
        ("left", ctypes.c_long),
        ("top", ctypes.c_long),
        ("right", ctypes.c_long),
        ("bottom", ctypes.c_long),
    ]


class MONITORINFOEXW(ctypes.Structure):
    _fields_ = [
        ("cbSize", ctypes.c_ulong),
        ("rcMonitor", RECT),
        ("rcWork", RECT),
        ("dwFlags", ctypes.c_ulong),
        ("szDevice", ctypes.c_wchar * 32),
    ]


class DISPLAY_DEVICEW(ctypes.Structure):
    _fields_ = [
        ("cb", ctypes.c_ulong),
        ("DeviceName", ctypes.c_wchar * 32),
        ("DeviceString", ctypes.c_wchar * 128),
        ("StateFlags", ctypes.c_ulong),
        ("DeviceID", ctypes.c_wchar * 128),
        ("DeviceKey", ctypes.c_wchar * 128),
    ]


MONITORENUMPROC = ctypes.WINFUNCTYPE(
    wintypes.BOOL,
    wintypes.HANDLE,
    wintypes.HDC,
    ctypes.POINTER(RECT),
    wintypes.LPARAM,
)


def list_monitors() -> list[MonitorInfo]:
    monitors: list[MonitorInfo] = []

    def callback(hmonitor, _hdc, _lprect, _lparam):
        info = MONITORINFOEXW()
        info.cbSize = ctypes.sizeof(MONITORINFOEXW)
        if not user32.GetMonitorInfoW(hmonitor, ctypes.byref(info)):
            return 1

        device = DISPLAY_DEVICEW()
        device.cb = ctypes.sizeof(DISPLAY_DEVICEW)
        friendly_name = info.szDevice
        if user32.EnumDisplayDevicesW(info.szDevice, 0, ctypes.byref(device), 0):
            friendly_name = device.DeviceString or info.szDevice

        width = info.rcMonitor.right - info.rcMonitor.left
        height = info.rcMonitor.bottom - info.rcMonitor.top
        label = f"{friendly_name} ({width}x{height}) [{info.rcMonitor.left}, {info.rcMonitor.top}]"
        monitors.append(
            MonitorInfo(
                label=label,
                device_name=info.szDevice,
                x=info.rcMonitor.left,
                y=info.rcMonitor.top,
                width=width,
                height=height,
            )
        )
        return 1

    user32.EnumDisplayMonitors(0, 0, MONITORENUMPROC(callback), 0)
    return monitors
