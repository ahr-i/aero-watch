@echo off
setlocal
cd /d "%~dp0"

set "PYTHON_CMD="
set "VENV_DIR=%~dp0.venv"
set "VENV_PYTHON=%VENV_DIR%\Scripts\python.exe"

where py >nul 2>nul
if %errorlevel%==0 (
    set "PYTHON_CMD=py"
) else (
    where python >nul 2>nul
    if %errorlevel%==0 (
        set "PYTHON_CMD=python"
    )
)

if not defined PYTHON_CMD (
    echo Python was not found. Install Python 3.10 or newer first.
    exit /b 1
)

if not exist "%VENV_PYTHON%" (
    echo Creating virtual environment...
    %PYTHON_CMD% -m venv "%VENV_DIR%"
    if errorlevel 1 (
        echo Failed to create virtual environment.
        exit /b 1
    )
)

echo Installing Python packages in .venv...
"%VENV_PYTHON%" -m pip install --upgrade pip
if errorlevel 1 (
    echo Failed to upgrade pip in virtual environment.
    exit /b 1
)

"%VENV_PYTHON%" -m pip install -r requirements.txt
if errorlevel 1 (
    echo Failed to install Python packages.
    exit /b 1
)

where ffmpeg >nul 2>nul
if %errorlevel%==0 (
    echo FFmpeg is already available on PATH.
    echo Install complete.
    exit /b 0
)

if exist "C:\ffmpeg\bin\ffmpeg.exe" (
    echo FFmpeg found at C:\ffmpeg\bin\ffmpeg.exe
    echo Install complete.
    exit /b 0
)

where winget >nul 2>nul
if not %errorlevel%==0 (
    echo FFmpeg was not found and winget is not available.
    echo Install FFmpeg manually, then update settings.json if needed.
    exit /b 1
)

echo FFmpeg was not found. Attempting installation with winget...
winget install --id Gyan.FFmpeg --accept-package-agreements --accept-source-agreements
if errorlevel 1 (
    echo FFmpeg installation failed. Install it manually, then update settings.json if needed.
    exit /b 1
)

echo Install complete.
