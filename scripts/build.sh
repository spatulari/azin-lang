# Build script for Azin compiler
# Usage: ./build.sh

AZC_SOURCE=cmd\compiler

go build -o azc.exe $AZC_SOURCE
