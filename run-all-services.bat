@echo off
REM Script to run all microservices in separate terminal windows

echo.
echo Starting all microservices...
echo ==============================
echo.

REM Change to the script directory
cd /d "%~dp0"

REM Start API Gateway
echo Starting API Gateway...
start "API Gateway" cmd /k "cd api-gateway && go run ./cmd/main.go"
timeout /t 2 /nobreak

REM Start User Service
echo Starting User Service...
start "User Service" cmd /k "cd user-service && go run ./cmd/main.go"
timeout /t 2 /nobreak

REM Start Inventory Service
echo Starting Inventory Service...
start "Inventory Service" cmd /k "cd inventory-service && go run ./cmd/main.go"
timeout /t 2 /nobreak

REM Start Order Service
echo Starting Order Service...
start "Order Service" cmd /k "cd order-service && go run ./cmd/main.go"
timeout /t 2 /nobreak

REM Start Mail Service
echo Starting Mail Service...
start "Mail Service" cmd /k "cd mail-service && go run ./cmd/main.go"
timeout /t 2 /nobreak

echo.
echo All services started in separate terminal windows
echo.
pause
