SHELL = /bin/sh

PREFIX = /usr/local
COMPLETIONS_DIR_BASH = $(PREFIX)/share/bash-completion/completions
COMPLETIONS_DIR_ZSH = $(PREFIX)/share/zsh/site-functions
COMPLETIONS_DIR_FISH = $(PREFIX)/share/fish/vendor_completions.d


.PHONY: all
all: miau completions

.PHONY: miau
miau:
	 go build -ldflags "-X main.version=$$(git describe --always --dirty)" .

.PHONY: completions
completions: miau.bash miau.zsh miau.fish

.PHONY: miau.bash
miau.bash: miau
	./miau --addr 127.0.0.1 completion bash > miau.bash

.PHONY: miau.zsh
miau.zsh: miau
	./miau --addr 127.0.0.1 completion zsh > miau.zsh

.PHONY: miau.fish
miau.fish: miau
	./miau --addr 127.0.0.1 completion fish > miau.fish

.PHONY: install
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

.PHONY: uninstall
uninstall:
	rm -f \
		$(PREFIX)/bin/miau \
		$(COMPLETIONS_DIR_BASH)/miau \
		$(COMPLETIONS_DIR_ZSH)/_miau \
		$(COMPLETIONS_DIR_FISH)/miau.fish

.PHONY: clean
clean: 
	rm -f miau miau.bash miau.zsh miau.fish
