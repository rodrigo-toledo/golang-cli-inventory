#!/bin/bash

echo "📊 Measuring test coverage..."
go test -coverprofile=coverage.out -covermode=count ./... > coverage_tmp.txt 2>&1

echo ""
echo "📈 Coverage by package:"
while IFS= read -r line; do
    if [[ $line == *"coverage:"* ]] && [[ $line != *"no statements"* ]]; then
        if [[ $line == ok* ]]; then
            # Format: ok   package/path   0.123s  coverage: 45.6% of statements
            package=$(echo "$line" | awk '{print $2}')
            coverage=$(echo "$line" | awk '{print $5}' | sed 's/%.*//')
            echo "$package: $coverage%"
        else
            # Format: package/path    coverage: 45.6% of statements
            package=$(echo "$line" | awk '{print $1}')
            coverage=$(echo "$line" | awk '{print $3}' | sed 's/%.*//')
            echo "$package: $coverage%"
        fi
    fi
done < coverage_tmp.txt | sort

echo ""
echo "📈 Overall coverage:"
go tool cover -func=coverage.out | grep "total:" | awk '{print $3}'

echo ""
echo "📄 HTML coverage report generated: coverage.html"
go tool cover -html=coverage.out -o coverage.html

rm -f coverage_tmp.txt