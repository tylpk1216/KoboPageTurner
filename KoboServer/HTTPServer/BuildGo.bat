@echo off

SET PATH=%PATH%;C:\Go\bin
SET GOPATH=%CD%

SET GOOS=linux
SET GOARCH=arm
SET GOARM=7

REM start cmd

SET EXE=HTTPServerARM

go build -o %EXE% HTTPServer.go

IF "%ERRORLEVEL%" == "1" (
    echo.
    echo Something is wrong.
    echo.
    pause
    exit
)

copy  /Y %EXE% ..\src\usr\local\koboserver

echo.
pause