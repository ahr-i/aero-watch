# streaming-test

Simple Windows Python app for:

- listing connected monitors
- showing a live preview of the selected monitor
- streaming the selected monitor to an RTMP URL with FFmpeg
- loading local app settings from `settings.json`

## Files

- `main.py`: app entry point
- `app/`: Python source package
- `app/ui.py`: Tkinter UI and app flow
- `app/config.py`: settings loading
- `app/monitoring.py`: Windows monitor lookup
- `app/ffmpeg.py`: FFmpeg discovery and stream command building
- `app/models.py`: shared data models
- `assets/icon/icon.png`: app icon
- `settings.json`: local settings such as FFmpeg path and default RTMP URL
- `install.bat`: creates `.venv`, installs Python packages there, and tries to install FFmpeg with winget
- `run.bat`: runs setup first, then starts the app with `.venv`

## Requirements

- Windows
- Python 3.10+
- FFmpeg installed, or `install.bat` can try to install it

Python packages:

```bash
install.bat
```

If FFmpeg is not on `PATH`, put its full path in `settings.json`:

```json
{
  "ffmpeg_path": "C:\\ffmpeg\\bin\\ffmpeg.exe"
}
```

## Install

```bash
install.bat
```

## Run

```bash
run.bat
```

## Notes

- The app currently streams video only. Audio is not included.
- RTMP URLs must start with `rtmp://`.
- The actual stream is sent with FFmpeg using Windows `gdigrab`.
