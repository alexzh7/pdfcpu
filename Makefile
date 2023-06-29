build: clean
	go build \
		-o bin/example_merge \
		cmd/example_merge/main.go

clean:
	@rm -rf bin/*
