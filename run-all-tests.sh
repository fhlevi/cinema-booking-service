#!/bin/bash

echo "========================================"
echo "Running Unit Tests for Cinema Booking System"
echo "========================================"

TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

echo ""
echo "[1/4] Testing Auth Service..."
echo "----------------------------------------"
cd auth-service
go mod tidy
if go test -v ./...; then
    echo "✅ Auth Service tests PASSED"
    ((PASSED_TESTS++))
else
    echo "❌ Auth Service tests FAILED"
    ((FAILED_TESTS++))
fi
((TOTAL_TESTS++))
cd ..

echo ""
echo "[2/4] Testing Cinema Service..."
echo "----------------------------------------"
cd cinema-service
go mod tidy
if go test -v ./...; then
    echo "✅ Cinema Service tests PASSED"
    ((PASSED_TESTS++))
else
    echo "❌ Cinema Service tests FAILED"
    ((FAILED_TESTS++))
fi
((TOTAL_TESTS++))
cd ..

echo ""
echo "[3/4] Testing Booking Service..."
echo "----------------------------------------"
cd booking-service
go mod tidy
if go test -v ./...; then
    echo "✅ Booking Service tests PASSED"
    ((PASSED_TESTS++))
else
    echo "❌ Booking Service tests FAILED"
    ((FAILED_TESTS++))
fi
((TOTAL_TESTS++))
cd ..

echo ""
echo "[4/4] Testing API Gateway..."
echo "----------------------------------------"
cd api-gateway
go mod tidy
if go test -v ./...; then
    echo "✅ API Gateway tests PASSED"
    ((PASSED_TESTS++))
else
    echo "❌ API Gateway tests FAILED"
    ((FAILED_TESTS++))
fi
((TOTAL_TESTS++))
cd ..

echo ""
echo "========================================"
echo "TEST SUMMARY"
echo "========================================"
echo "Total Services: $TOTAL_TESTS"
echo "Passed: $PASSED_TESTS"
echo "Failed: $FAILED_TESTS"

if [ $FAILED_TESTS -eq 0 ]; then
    echo ""
    echo "🎉 ALL TESTS PASSED! 🎉"
    exit 0
else
    echo ""
    echo "⚠️  Some tests failed. Please check the output above."
    exit 1
fi