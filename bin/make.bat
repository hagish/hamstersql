cd ..
cd hamstersql

# linux binary
set GOARCH=amd64
set GOOS=linux
go build

# windows binary
set GOARCH=amd64
set GOOS=windows
go build

move /Y hamstersql.exe ..\bin\hamstersql.exe
move /Y hamstersql ..\bin\hamstersql

cd ..
cd bin

