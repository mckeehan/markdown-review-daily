.PHONY: build install clean run test

BINARY_NAME=markdown_review_daily
INSTALL_PATH=/usr/local/bin

build: deps $(BINARY_NAME)

install: build
	@echo "Installing to $(INSTALL_PATH)..."
	sudo mv $(BINARY_NAME) $(INSTALL_PATH)/
	@echo "Installed successfully!"

clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)
	@echo "Clean complete!"

run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BINARY_NAME)

test:
	@echo "Running tests..."
	go test -v ./...

deps: go.sum

go.sum: go.mod
	@echo "Downloading dependencies..."
	go mod tidy
	@echo "Dependencies downloaded!"

config:
	@if [ ! -f config.yaml ]; then \
		echo "Creating default config.yaml..."; \
		cp config.yaml.example config.yaml; \
		echo "Please edit config.yaml with your settings"; \
	else \
		echo "config.yaml already exists"; \
	fi

# markdown_review_daily: markdown_review_daily.go
# 	@echo "Building $(BINARY_NAME)..."
# 	go build -o $(BINARY_NAME) markdown_review_daily.go
# 	@echo "Build complete!"

%: %.go
	@echo "Building $@..."
	go build -o $@ $<
	@echo "Build complete!"
