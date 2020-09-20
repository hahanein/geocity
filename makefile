all: dist/main.wasm dist/main

dist/main.wasm: ./frontend/*
	GOOS=js GOARCH=wasm go build -o ./dist/main.wasm ./frontend/

dist/main: ./backend/*
	go build -o ./dist/main ./backend/

clean:
	rm -rf ./dist/*