#!/bin/bash

# debug on/off (-/+)
set +x

# destination
filename="osmnotes2csv"

echo "Building Linux ..."
env GOOS=linux GOARCH=amd64 go build
zip "$filename"_linux_amd64.zip "$filename"
rm "$filename"

echo "Building MacOS ..."
env GOOS=darwin GOARCH=amd64 go build
zip "$filename"_darwin_amd64.zip "$filename"
rm "$filename"

echo "Building Windows ..."
env GOOS=windows GOARCH=amd64 go build
zip "$filename"_windows_amd64.zip "$filename".exe
rm "$filename".exe
