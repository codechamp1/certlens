# Variables
BINARY_NAME=certlens
PKG=./...
GO=go

# Creates kubectl plugin
.PHONY: kubectl-plugin
kubectl-plugin:
	@echo "Installing kubectl-certlens plugin..."
	@command -v certlens >/dev/null 2>&1 || { echo "error: certlens not found in PATH."; exit 1; }
	@command -v kubectl >/dev/null 2>&1 || { echo "error: kubectl not found in PATH."; exit 1; }
	@mkdir -p "$$HOME/.local/bin"
	@ln -sf "$$(command -v certlens)" "$$HOME/.local/bin/kubectl-certlens"
	@chmod +x "$$HOME/.local/bin/kubectl-certlens"
	@echo "Installation complete!"

# Build executable
build:
	$(GO) build -o $(BINARY_NAME) ./cmd/$(BINARY_NAME)/

# run tests
test:
	$(GO) test -v -race -coverprofile=coverage.out $(PKG)

# test coverage
cover: test
	$(GO) tool cover -html=coverage.out


# run linter
lint:
	golangci-lint run

#  clean coverage
clean:
	rm -f $(BINARY_NAME) coverage.out

# run locally
run: build
	./$(BINARY_NAME)

# Install in $GOPATH/bin (or $GOBIN)
install:
	$(GO) install ./cmd/$(BINARY_NAME)/

action-ci: test
	@echo " ✅ All checks passed "


.PHONY: build test cover lint clean run install
