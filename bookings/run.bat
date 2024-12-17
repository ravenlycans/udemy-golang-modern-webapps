@echo off

:start
echo Building and Running Go Project
go build -o bookings.exe .\cmd\web
echo Build Complete.
echo Running the built exe.
bookings.exe

echo Program ended, deleting bookings.exe
del -F -Q bookings.exe
echo Done, rerunning.
goto start
