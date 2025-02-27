# Go variables
GO = go
MAIN_DIR = main
MAIN_FILE = $(MAIN_DIR)/main.go
BINARY_NAME = main


all: run

run:
	$(GO) run $(MAIN_FILE)
