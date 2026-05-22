import glob
import os
import shutil

from app.models import AppSettings, MonitorInfo


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


def build_stream_command(
    ffmpeg_path: str,
    monitor: MonitorInfo,
    url: str,
    settings: AppSettings,
) -> list[str]:
    bitrate = f"{settings.video_bitrate_kbps}k"
    bufsize = f"{settings.video_bitrate_kbps * 2}k"
    gop = str(settings.stream_fps * 2)

    return [
        ffmpeg_path,
        "-y",
        "-thread_queue_size",
        "512",
        "-f",
        "gdigrab",
        "-draw_mouse",
        "1",
        "-framerate",
        str(settings.stream_fps),
        "-offset_x",
        str(monitor.x),
        "-offset_y",
        str(monitor.y),
        "-video_size",
        f"{monitor.width}x{monitor.height}",
        "-i",
        "desktop",
        "-an",
        "-vsync",
        "cfr",
        "-c:v",
        "libx264",
        "-preset",
        "ultrafast",
        "-tune",
        "zerolatency",
        "-pix_fmt",
        "yuv420p",
        "-bf",
        "0",
        "-g",
        gop,
        "-b:v",
        bitrate,
        "-maxrate",
        bitrate,
        "-bufsize",
        bufsize,
        "-f",
        "flv",
        url,
    ]
