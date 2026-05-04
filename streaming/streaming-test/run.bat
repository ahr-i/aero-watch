@echo off
setlocal
cd /d "%~dp0"

set "VENV_PYTHON=%~dp0.venv\Scripts\python.exe"

call install.bat
if errorlevel 1 (
    echo Setup failed. Fix the issue above, then try again.
    exit /b 1
)

if not exist "%VENV_PYTHON%" (
    echo Virtual environment Python was not found.
    exit /b 1
)

"%VENV_PYTHON%" main.py
