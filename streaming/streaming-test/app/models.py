from dataclasses import dataclass


@dataclass
class MonitorInfo:
    label: str
    device_name: str
    x: int
    y: int
    width: int
    height: int


@dataclass
class AppSettings:
    ffmpeg_path: str
    default_rtmp_url: str
    preview_interval_ms: int
