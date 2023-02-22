install:
	go install .
	go run . --addr 127.0.0.1 completion fish > ~/.config/fish/completions/miau.fish
