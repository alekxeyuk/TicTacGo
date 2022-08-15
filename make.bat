@echo off
if %1==build (go build -o TicTacGo.exe -ldflags "-w -s")
if %1==release (git push --atomic origin master %2)
if %1==tag (git tag -a %2 -m "New release: %2")