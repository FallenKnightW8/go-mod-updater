@echo off
setlocal enabledelayedexpansion

if not exist "repos.txt" (
    echo Error: repos.txt not found!
    echo Creating example repos.txt...
    (
        echo gin-gonic/gin
        echo gorilla/mux
        echo labstack/echo
        echo gofiber/fiber
    ) > repos.txt
    echo Please edit repos.txt and run again.
    pause
    exit /b 1
)

:MAIN_MENU
cls
echo ======================================
echo Available repositories:
echo ======================================
set COUNT=0
for /f "tokens=* delims=" %%a in (repos.txt) do (
    set /a COUNT+=1
    echo [!COUNT!] %%a
)

if %COUNT% equ 0 (
    echo Error: repos.txt is empty!
    pause
    exit /b 1
)

echo ======================================
echo [0] Exit
echo ======================================
set /p CHOICE="Choose number (1-%COUNT%) or 0 to exit: "

if "%CHOICE%"=="0" exit /b 0

set VALID=0
for /l %%i in (1,1,%COUNT%) do (
    if "%CHOICE%"=="%%i" set VALID=1
)

if %VALID% equ 0 (
    echo Invalid choice! Please enter number between 1 and %COUNT%.
    pause
    goto MAIN_MENU
)

set NUM=0
for /f "tokens=* delims=" %%a in (repos.txt) do (
    set /a NUM+=1
    if !NUM! equ %CHOICE% set "REPO=%%a"
)

:FORMAT_MENU
cls
echo ======================================
echo Selected repository: %REPO%
echo ======================================
echo Choose output format:
echo [1] Save as JSON file
echo [ENTER] Show table and return to menu
echo ======================================
set /p FORMAT_CHOICE="Enter choice: "

if "%FORMAT_CHOICE%"=="1" goto SAVE_JSON

:: Если нажат Enter или что-то другое - стандартный вывод
cls
echo ======================================
echo Analyzing: %REPO%
echo ======================================
echo.

echo %REPO% | findstr /c:"http" >nul
if %errorlevel% equ 0 (
    go-mod-updater.exe --repo %REPO%
) else (
    go-mod-updater.exe --repo https://github.com/%REPO%
)

echo.
echo ======================================
echo Analysis complete!
echo ======================================
pause
goto MAIN_MENU

:SAVE_JSON
cls
echo ======================================
echo Analyzing and saving to JSON...
echo ======================================
echo.

:: Создаем папку results если её нет
if not exist "results" mkdir results

:: Формируем имя файла из репозитория и даты
for /f "tokens=1-4 delims=/. " %%a in ('date /t') do set MYDATE=%%a-%%b-%%c
for /f "tokens=1-2 delims=: " %%a in ('time /t') do set MYTIME=%%a-%%b

:: Заменяем слеши в имени репозитория на подчеркивания
set "FILENAME=%REPO:/=_%"
set "FILENAME=%FILENAME: =_%"
set "FILENAME=results\%FILENAME%_%MYDATE%_%MYTIME%.json"

echo Repository: %REPO%
echo Saving to: %FILENAME%
echo.

echo %REPO% | findstr /c:"http" >nul
if %errorlevel% equ 0 (
    go-mod-updater.exe --repo %REPO% --json > "%FILENAME%" 2>&1
) else (
    go-mod-updater.exe --repo https://github.com/%REPO% --json > "%FILENAME%" 2>&1
)

if %errorlevel% equ 0 (
    echo.
    echo ======================================
    echo JSON saved successfully!
    echo File: %FILENAME%
    echo ======================================
) else (
    echo.
    echo ======================================
    echo Error during analysis!
    echo Check file for details: %FILENAME%
    echo ======================================
)

pause
goto MAIN_MENU