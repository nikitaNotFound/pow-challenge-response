.PHONY: all build clean run-server run-client build-w clean-w run-server-w run-client-w run-all run-all-w

# Unix commands
BINARY_DIR = bin
SERVER_BINARY = $(BINARY_DIR)/server
CLIENT_BINARY = $(BINARY_DIR)/client

build:
	@mkdir -p $(BINARY_DIR)
	@echo "Building server..."
	@go build -o $(SERVER_BINARY) ./cmd/server
	@echo "Building client..."
	@go build -o $(CLIENT_BINARY) ./cmd/client

clean:
	@echo "Cleaning..."
	@rm -rf $(BINARY_DIR)

run-server: build
	@echo "Starting server..."
	@./$(SERVER_BINARY)

run-client: build
	@echo "Starting client..."
	@./$(CLIENT_BINARY)

run-all: build
	@echo "Starting server in background..."
	@./$(SERVER_BINARY) & \
	echo "Starting client..." && \
	./$(CLIENT_BINARY)

run-monitoring:
	expvarmon -ports="1234" -vars="mem:memstats.Alloc,mem:memstats.Sys,mem:memstats.HeapAlloc,mem:memstats.HeapInuse,duration:memstats.PauseNs,duration:memstats.PauseTotalNs,Goroutines,duration:Uptime,duration:Response.Mean"

# Windows commands
WIN_BINARY_DIR = bin
WIN_SERVER = $(WIN_BINARY_DIR)\server.exe
WIN_CLIENT = $(WIN_BINARY_DIR)\client.exe

build-w:
	@if not exist $(WIN_BINARY_DIR) mkdir $(WIN_BINARY_DIR)
	@echo "Building server..."
	@go build -o $(WIN_SERVER) ./cmd/server
	@echo "Building client..."
	@go build -o $(WIN_CLIENT) ./cmd/client

clean-w:
	@echo "Cleaning..."
	@if exist $(WIN_BINARY_DIR) rd /s /q $(WIN_BINARY_DIR)

run-server-w: build-w
	@echo "Starting server..."
	@.\$(WIN_SERVER)

run-client-w: build-w
	@echo "Starting client..."
	@.\$(WIN_CLIENT)

run-all-w: build-w
	@echo "Starting server in background..."
	@taskkill /F /IM server.exe >NUL 2>&1 || exit 0
	@timeout 1 >NUL
	@start /B cmd /c ".\$(WIN_SERVER)"
	@timeout 2 >NUL
	@echo "Starting client..."
	@.\$(WIN_CLIENT)

# Default to Windows if on Windows, Unix otherwise
ifeq ($(OS),Windows_NT)
all: build-w
else
all: build
endif 