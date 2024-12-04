# Define the name of the Go binary
BINARY_NAME=corb4n-c2

# Define the directory where your frontend files are located
FRONTEND_DIR=frontend

# Define the build command for the frontend (e.g., using npm or any other build tool)
FRONTEND_BUILD_CMD=cd $(FRONTEND_DIR) && npm run build

# Define the build command for the Go server
GO_BUILD_CMD=go build -o $(BINARY_NAME) .

# Define the run command for the Go server
RUN_CMD=go run .

# Define the default target
all: build-frontend run

# Define the build-frontend target
build-frontend:
	$(FRONTEND_BUILD_CMD)

# Define the build-go target
build-go:
	$(GO_BUILD_CMD)

# Define the run target
run: build-frontend
	$(RUN_CMD)

# Define the build target that builds both frontend and Go executable
build: build-frontend build-go

# Define the clean target
clean:
	rm -f $(BINARY_NAME)

.PHONY: all build-frontend build-go run build clean
