# Makefile
APP_NAME=bilidata

build:
    go build -o $(APP_NAME)

run: build
    ./$(APP_NAME)

clean:
    rm -f $(APP_NAME)

