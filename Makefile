
all: build

# Set OS-specific variables
ifeq ($(OS),Windows_NT)
    TARGET = output/medal_win.exe
else
    ifeq ($(shell uname -s),Linux)
        TARGET = output/medal_linux
    endif
endif

build:
	go build -o $(TARGET) -ldflags '-w -s' main.go

run: build
	$(TARGET)

.PHONY: all build run
