mkdir -p build
commitdate=$(git show -s --format=%ci | sed 's/[: -]/_/g' | sed 's/_+.*//g')
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -ldflags="-s -w" -o tankzors_$commitdate.exe *.go
echo "Archive base name?"
read basename
zip build/$basename$commitdate.zip -r tankzors_$commitdate.exe assets/*
rm tankzors_$commitdate.exe
echo "Build successful."