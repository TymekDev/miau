PREFIX = /usr/local
COMPLETIONS_DIR_BASH = $(PREFIX)/share/bash-completion/completions
COMPLETIONS_DIR_ZSH = $(PREFIX)/share/zsh/site-functions
COMPLETIONS_DIR_FISH = $(PREFIX)/share/fish/vendor_completions.d

all: miau completions

miau:
	 go build -ldflags "-X main.version=$$(git describe --always --dirty)" .

completions: miau.bash miau.zsh miau.fish

miau.bash: miau
	./miau --addr 127.0.0.1 completion bash > miau.bash

miau.zsh: miau
	./miau --addr 127.0.0.1 completion zsh > miau.zsh

miau.fish: miau
	./miau --addr 127.0.0.1 completion fish > miau.fish

clean: 
	rm -f miau miau.bash miau.zsh miau.fish

install:
	install -d \
		$(PREFIX)/bin \
		$(COMPLETIONS_DIR_BASH) \
		$(COMPLETIONS_DIR_ZSH) \
		$(COMPLETIONS_DIR_FISH)

	install -pm 0755 miau $(PREFIX)/bin/miau
	install -pm 0644 miau.bash $(COMPLETIONS_DIR_BASH)/miau
	install -pm 0644 miau.zsh $(COMPLETIONS_DIR_ZSH)/_miau
	install -pm 0644 miau.fish $(COMPLETIONS_DIR_FISH)/miau.fish

uninstall:
	rm -f \
		$(PREFIX)/bin/miau \
		$(COMPLETIONS_DIR_BASH)/miau \
		$(COMPLETIONS_DIR_ZSH)/_miau \
		$(COMPLETIONS_DIR_FISH)/miau.fish

.PHONY: all miau completions clean install uninstall
