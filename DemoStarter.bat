@echo off
set /p slavecount="how many slaves: "
set /p artificialDelay="artificial network delay(ms): "
set /P artificialOffsetMin="artificial offset minimum(ms): "
set /P artificialOffsetMax="artificial offset maximum(ms): "
set /p verbose="Verbose output? (true/false): "
SET /A test=%RANDOM% * 100 / 32768 + 1

setlocal ENABLEDELAYEDEXPANSION
start cmd.exe /k go run master\master.go 5010 6666 %artificialDelay% %verbose%
set /a port =6060
SET /A offset=!RANDOM!%%artificialOffsetMax+ artificialOffsetMin 
for /L %%x in (1, 1, %slavecount%) do (
	SET /A offset=!RANDOM!%%artificialOffsetMax + artificialOffsetMin
	set /a port+=1
	start cmd.exe /k go run slave\slave.go !port! 127.0.0.1:5010 224.0.1.1:6666 !offset! %artificialDelay% %verbose%
)