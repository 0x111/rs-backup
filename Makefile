.PHONY: all
all: build_linux_amd64 build_darwin_amd64 build_windows_amd64 checksums

.PHONY: build_linux_amd64
build_linux_amd64:
	GOOS=linux GOARCH=amd64 go build -v -a -gcflags=-trimpath=$$PWD -asmflags=-trimpath=$$PWD -o build/rs-backup-linux-amd64

.PHONY: build_linux_i386
build_linux_i386:
	GOOS=linux GOARCH=386 go build -v -a -gcflags=-trimpath=$$PWD -asmflags=-trimpath=$$PWD -o build/rs-backup-linux-i386

.PHONY: build_darwin_amd64
build_darwin_amd64:
	GOOS=darwin GOARCH=amd64 go build -v -a -gcflags=-trimpath=$$PWD -asmflags=-trimpath=$$PWD -o build/rs-backup-darwin-amd64

.PHONY: build_darwin_i386
build_darwin_i386:
	GOOS=darwin GOARCH=386 go build -v -a -gcflags=-trimpath=$$PWD -asmflags=-trimpath=$$PWD -o build/rs-backup-darwin-i386

.PHONY: build_windows_amd64
build_windows_amd64:
	CC=/usr/local/bin/x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build -v -a -gcflags=-trimpath=$$PWD -asmflags=-trimpath=$$PWD -o build/rs-backup-windows-amd64.exe

.PHONY: build_windows_i386
build_windows_i386:
	CC=/usr/local/bin/x86_64-w64-mingw32-gcc GOOS=windows GOARCH=386 go build -v -a -gcflags=-trimpath=$$PWD -asmflags=-trimpath=$$PWD -o build/rs-backup-windows-i386.exe

.PHONY: checksums
checksums:
	shasum -a 256 build/* > build/checksum.txt

test: