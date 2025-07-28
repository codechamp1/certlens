# Variabile
BINARY_NAME=certlens
PKG=./...
GO=go

# Build executabilul
build:
	$(GO) build -o $(BINARY_NAME) ./cmd/$(BINARY_NAME)/

# Rulează testele cu raport de acoperire
test:
	$(GO) test -v -race -coverprofile=coverage.out $(PKG)

# Raport de acoperire (deschide în browser)
cover: test
	$(GO) tool cover -html=coverage.out


# Rulează linterul golangci-lint (trebuie instalat golangci-lint)
lint: install-lint
	golangci-lint run

# Curățenie - șterge binarele și fișierele generate
clean:
	rm -f $(BINARY_NAME) coverage.out

# Rulează aplicația local (depinde de build)
run: build
	./$(BINARY_NAME)

# Instalează binarul în $GOPATH/bin (sau $GOBIN)
install:
	$(GO) install ./cmd/$(BINARY_NAME)/

action: lint test
	@echo " ✅ All checks passed "


.PHONY: build test cover lint clean run install
