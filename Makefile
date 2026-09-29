# koolkit dev Makefile. `make install.dev` builds the kk dispatcher onto
# PATH and every sibling plugin into the managed plugin dir, so local
# source changes run through the real `kk <plugin>` dispatch path.

XDG_DATA_HOME ?= $(HOME)/.local/share
MANAGED_BIN := $(XDG_DATA_HOME)/koolkit/bin
GOBIN := $(shell go env GOPATH)/bin
PLUGINS := kk-goto kk-prs

.PHONY: install.dev install.dev.kk install.dev.plugins

install.dev: install.dev.kk install.dev.plugins

install.dev.kk:
	@mkdir -p "$(GOBIN)"
	go build -o "$(GOBIN)/kk" .
	@echo "installed kk -> $(GOBIN)/kk"

install.dev.plugins:
	@mkdir -p "$(MANAGED_BIN)"
	@for p in $(PLUGINS); do \
		if [ -d "../$$p" ]; then \
			echo "building $$p -> $(MANAGED_BIN)/$$p"; \
			( cd "../$$p" && go build -o "$(MANAGED_BIN)/$$p" . ); \
		else \
			echo "skip $$p (../$$p not found)"; \
		fi; \
	done
