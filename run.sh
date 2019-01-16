#! /bin/bash

RPATH=`realpath $0`
DIR=`dirname $RPATH`

echo Checking can package exists 
go get github.com/brutella/can

echo Checking sqlite3 package exists 
go get github.com/mattn/go-sqlite3

go run $DIR/main/*.go
