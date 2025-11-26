@echo off
echo ========================================
echo Generating Test Coverage Reports
echo ========================================

echo.
echo [1/4] Auth Service Coverage...
echo ----------------------------------------
cd auth-service
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
echo Coverage report generated: auth-service/coverage.html
cd ..

echo.
echo [2/4] Cinema Service Coverage...
echo ----------------------------------------
cd cinema-service
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
echo Coverage report generated: cinema-service/coverage.html
cd ..

echo.
echo [3/4] Booking Service Coverage...
echo ----------------------------------------
cd booking-service
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
echo Coverage report generated: booking-service/coverage.html
cd ..

echo.
echo [4/4] API Gateway Coverage...
echo ----------------------------------------
cd api-gateway
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
echo Coverage report generated: api-gateway/coverage.html
cd ..

echo.
echo ========================================
echo Coverage reports generated successfully!
echo Open the .html files in your browser to view detailed coverage.
echo ========================================