build: clean
	mkdir -p bin
	go build -o bin/danglingTpls

clean:
	rm -rf bin
