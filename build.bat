@echo off
setlocal EnableDelayedExpansion
title SimRacing Dashboard - Build

echo.
echo +========================================+
echo ^|   SimRacing Dashboard - Build          ^|
echo +========================================+
echo.

set SCRIPT_DIR=%~dp0
set FRONTEND_DIR=%SCRIPT_DIR%frontend
set BACKEND_DIR=%SCRIPT_DIR%backend

:: Step 1: Frontend dependencies
echo ^> Step 1: Installing Nuxt dependencies...
cd /d "%FRONTEND_DIR%"
call npm install
if %ERRORLEVEL% neq 0 ( echo ERROR: npm install failed & exit /b 1 )
echo   Done.

:: Step 2: Build Nuxt
echo.
echo ^> Step 2: Building Nuxt frontend...
call npx nuxi generate
if %ERRORLEVEL% neq 0 ( echo ERROR: nuxi generate failed & exit /b 1 )
echo   Done. Output to backend\static\

if not exist "%BACKEND_DIR%\static\index.html" (
  echo ERROR: index.html not found - build failed
  exit /b 1
)

:: Step 3: Go dependencies
echo.
echo ^> Step 3: Downloading Go modules...
cd /d "%BACKEND_DIR%"
go mod download
if %ERRORLEVEL% neq 0 ( echo ERROR: go mod download failed & exit /b 1 )
echo   Done.

:: Step 4: Build EXE
echo.
echo ^> Step 4: Compiling Windows EXE...
set CGO_ENABLED=0
set GOOS=windows
set GOARCH=amd64
go build -ldflags="-s -w" -o "%SCRIPT_DIR%simracing-dashboard.exe" .
if %ERRORLEVEL% neq 0 ( echo ERROR: go build failed & exit /b 1 )
echo   Done.

echo.
echo +========================================+
echo ^|  Build complete!                       ^|
echo +========================================+
echo.
echo   Run: simracing-dashboard.exe
echo   Simulation: simracing-dashboard.exe --simulate
echo   Custom port: simracing-dashboard.exe --addr :9090
echo.
pause
