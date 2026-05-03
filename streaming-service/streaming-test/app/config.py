import json
import os

from app.models import AppSettings


APP_DIR = os.path.dirname(os.path.abspath(__file__))
BASE_DIR = os.path.dirname(APP_DIR)
SETTINGS_PATH = os.path.join(BASE_DIR, "settings.json")


def default_settings() -> AppSettings:
    return AppSettings(
        ffmpeg_path="",
        default_rtmp_url="rtmp://",
        preview_interval_ms=700,
    )


def load_settings() -> AppSettings:
    settings = default_settings()
    if not os.path.exists(SETTINGS_PATH):
        return settings

    try:
        with open(SETTINGS_PATH, "r", encoding="utf-8") as file:
            loaded = json.load(file)
    except (OSError, json.JSONDecodeError):
        return settings

    if not isinstance(loaded, dict):
        return settings

    return AppSettings(
        ffmpeg_path=str(loaded.get("ffmpeg_path", settings.ffmpeg_path)),
        default_rtmp_url=str(loaded.get("default_rtmp_url", settings.default_rtmp_url)),
        preview_interval_ms=_get_preview_interval_ms(loaded.get("preview_interval_ms")),
    )


def _get_preview_interval_ms(value: object) -> int:
    try:
        interval_ms = int(value)
    except (TypeError, ValueError):
        return 700
    return max(100, interval_ms)
