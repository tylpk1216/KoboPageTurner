@echo off

SET PATH=%PATH%;C:\Go\bin
SET GOPATH=%CD%

REM start cmd

SET EXE=Win.exe

go build -o %EXE% Win.go

IF NOT "%ERRORLEVEL%" == "0" (
    echo.
    echo Something is wrong.
    echo.
    pause
    exit
)

echo.
pause