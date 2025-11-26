@echo off
echo ========================================
echo Running Unit Tests for Cinema Booking System
echo ========================================

set TOTAL_TESTS=0
set PASSED_TESTS=0
set FAILED_TESTS=0

echo.
echo [1/4] Testing Auth Service...
echo ----------------------------------------
cd auth-service
go mod tidy
go test -v ./...
if %ERRORLEVEL% EQU 0 (
    echo ✅ Auth Service tests PASSED
    set /a PASSED_TESTS+=1
) else (
    echo ❌ Auth Service tests FAILED
    set /a FAILED_TESTS+=1
)
set /a TOTAL_TESTS+=1
cd ..

echo.
echo [2/4] Testing Cinema Service...
echo ----------------------------------------
cd cinema-service
go mod tidy
go test -v ./...
if %ERRORLEVEL% EQU 0 (
    echo ✅ Cinema Service tests PASSED
    set /a PASSED_TESTS+=1
) else (
    echo ❌ Cinema Service tests FAILED
    set /a FAILED_TESTS+=1
)
set /a TOTAL_TESTS+=1
cd ..

echo.
echo [3/4] Testing Booking Service...
echo ----------------------------------------
cd booking-service
go mod tidy
go test -v ./...
if %ERRORLEVEL% EQU 0 (
    echo ✅ Booking Service tests PASSED
    set /a PASSED_TESTS+=1
) else (
    echo ❌ Booking Service tests FAILED
    set /a FAILED_TESTS+=1
)
set /a TOTAL_TESTS+=1
cd ..

echo.
echo [4/4] Testing API Gateway...
echo ----------------------------------------
cd api-gateway
go mod tidy
go test -v ./...
if %ERRORLEVEL% EQU 0 (
    echo ✅ API Gateway tests PASSED
    set /a PASSED_TESTS+=1
) else (
    echo ❌ API Gateway tests FAILED
    set /a FAILED_TESTS+=1
)
set /a TOTAL_TESTS+=1
cd ..

echo.
echo ========================================
echo TEST SUMMARY
echo ========================================
echo Total Services: %TOTAL_TESTS%
echo Passed: %PASSED_TESTS%
echo Failed: %FAILED_TESTS%

if %FAILED_TESTS% EQU 0 (
    echo.
    echo 🎉 ALL TESTS PASSED! 🎉
    exit /b 0
) else (
    echo.
    echo ⚠️  Some tests failed. Please check the output above.
    exit /b 1
)