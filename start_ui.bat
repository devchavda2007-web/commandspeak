@echo off
echo Starting CommandSpeak Frontend...

:: Check if PHP is installed
php -v >nul 2>&1
IF %ERRORLEVEL% NEQ 0 (
    echo [WARNING] PHP is not installed or not in PATH!
    echo Falling back to static mode (localStorage).
    echo Opening index.html directly...
    start index.html
    exit /b
)

:: If PHP is installed, start the local server
echo PHP is installed. Starting private SQLite backend on localhost:8080...
echo (Keep this window open to keep the server running)
echo.

:: Open the browser
start http://localhost:8080/index.html

:: Run the server
php -S localhost:8080
