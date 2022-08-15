@echo off
if %1==build (go build -o TicTacGo.exe -ldflags "-w -s")
if %1==release (git push --atomic origin master %2)