#!/bin/bash

# LinkSprint Load Testing Script
# This script tests the application performance using curl

echo "🚀 Starting LinkSprint Load Test"
echo "=================================="

# Configuration
BASE_URL="http://localhost:8080"
REQUESTS=1000
CONCURRENT=10

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test 1: Health Check
echo -e "${YELLOW}Testing Health Check...${NC}"
start_time=$(date +%s.%N)
for i in {1..100}; do
    curl -s "$BASE_URL/health" > /dev/null
done
end_time=$(date +%s.%N)
health_time=$(echo "$end_time - $start_time" | bc)
echo -e "${GREEN}Health check completed in ${health_time}s${NC}"

# Test 2: URL Creation
echo -e "${YELLOW}Testing URL Creation...${NC}"
start_time=$(date +%s.%N)
for i in {1..50}; do
    curl -s -X POST "$BASE_URL/api/v1/urls/shorten" \
        -H "Content-Type: application/json" \
        -d "{\"original_url\":\"https://example$i.com\"}" > /dev/null
done
end_time=$(date +%s.%N)
create_time=$(echo "$end_time - $start_time" | bc)
echo -e "${GREEN}URL creation completed in ${create_time}s${NC}"

# Test 3: URL Listing
echo -e "${YELLOW}Testing URL Listing...${NC}"
start_time=$(date +%s.%N)
for i in {1..20}; do
    curl -s "$BASE_URL/api/v1/urls" > /dev/null
done
end_time=$(date +%s.%N)
list_time=$(echo "$end_time - $start_time" | bc)
echo -e "${GREEN}URL listing completed in ${list_time}s${NC}"

# Test 4: Analytics
echo -e "${YELLOW}Testing Analytics...${NC}"
start_time=$(date +%s.%N)
for i in {1..30}; do
    curl -s "$BASE_URL/api/v1/analytics/global" > /dev/null
done
end_time=$(date +%s.%N)
analytics_time=$(echo "$end_time - $start_time" | bc)
echo -e "${GREEN}Analytics completed in ${analytics_time}s${NC}"

# Summary
echo ""
echo -e "${GREEN}Load Test Summary:${NC}"
echo "========================"
echo -e "Health Checks: ${GREEN}100 requests${NC}"
echo -e "URL Creation: ${GREEN}50 requests${NC}"
echo -e "URL Listing:  ${GREEN}20 requests${NC}"
echo -e "Analytics:    ${GREEN}30 requests${NC}"
echo ""
echo -e "Total Time: ${YELLOW}${health_time}s + ${create_time}s + ${list_time}s + ${analytics_time}s${NC}"

# Performance metrics
total_requests=200
total_time=$(echo "$health_time + $create_time + $list_time + $analytics_time" | bc)
avg_response_time=$(echo "scale=3; $total_time / $total_requests" | bc)
requests_per_second=$(echo "scale=2; $total_requests / $total_time" | bc)

echo ""
echo -e "${GREEN}Performance Metrics:${NC}"
echo "========================"
echo -e "Total Requests: ${GREEN}$total_requests${NC}"
echo -e "Total Time: ${GREEN}${total_time}s${NC}"
echo -e "Average Response Time: ${GREEN}${avg_response_time}s${NC}"
echo -e "Requests per Second: ${GREEN}${requests_per_second}${NC}"

echo ""
echo -e "${GREEN}Load test completed! 🎉${NC}" 