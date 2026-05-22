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
    stream_fps: int
    video_bitrate_kbps: int
    pause_preview_during_stream: bool
