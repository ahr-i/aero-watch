import glob
import os
import shutil

from app.models import MonitorInfo


def find_ffmpeg(configured_path: str = "") -> str | None:
    if configured_path and os.path.exists(configured_path):
        return configured_path

    env_path = os.environ.get("AEROWATCH_FFMPEG")
    if env_path and os.path.exists(env_path):
        return env_path

    on_path = shutil.which("ffmpeg")
    if on_path:
        return on_path

    search_roots = [
        os.path.expandvars(
            r"%LOCALAPPDATA%\Microsoft\WinGet\Packages\Gyan.FFmpeg_Microsoft.Winget.Source_8wekyb3d8bbwe"
        ),
        r"C:\ffmpeg",
    ]
    for root in search_roots:
        pattern = os.path.join(root, "**", "ffmpeg.exe")
        matches = glob.glob(pattern, recursive=True)
        if matches:
            return matches[0]

    return None


def build_stream_command(ffmpeg_path: str, monitor: MonitorInfo, url: str) -> list[str]:
    return [
        ffmpeg_path,
        "-y",
        "-f",
        "gdigrab",
        "-framerate",
        "30",
        "-offset_x",
        str(monitor.x),
        "-offset_y",
        str(monitor.y),
        "-video_size",
        f"{monitor.width}x{monitor.height}",
        "-i",
        "desktop",
        "-an",
        "-c:v",
        "libx264",
        "-preset",
        "veryfast",
        "-pix_fmt",
        "yuv420p",
        "-g",
        "60",
        "-b:v",
        "6000k",
        "-maxrate",
        "6000k",
        "-bufsize",
        "12000k",
        "-f",
        "flv",
        url,
    ]
