@echo off

SET PATH=%PATH%;C:\Go\bin
SET GOPATH=%CD%

SET GOOS=linux
SET GOARCH=arm
SET GOARM=7

REM start cmd

SET EXE=NewDeviceDebug

go build -o %EXE% NewDevice.go

IF NOT "%ERRORLEVEL%" == "0" (
    echo.
    echo Something is wrong.
    echo.
    pause
    exit
)

copy  /Y %EXE% ..\src\usr\local\koboserver

echo.
pause