@echo off
echo Initializing...
echo.
if "%~1" == "" (
	src\decode.exe
) else (
	src\decode.exe "%~1"
)

echo.
pause