#!/bin/bash

# Network Traffic Generator Script
# This script generates various types of network traffic for testing

echo "🚀 Starting network traffic generation for PCAP testing..."

# Function to generate HTTP traffic
generate_http_traffic() {
    echo "📡 Generating HTTP traffic..."
    for i in {1..10}; do
        curl -s https://httpbin.org/get?param=$i > /dev/null &
        curl -s https://jsonplaceholder.typicode.com/posts/$i > /dev/null &
        sleep 0.5
    done
    wait
    echo "✅ HTTP traffic generated"
}

# Function to generate DNS queries
generate_dns_traffic() {
    echo "🌐 Generating DNS traffic..."
    domains=("google.com" "github.com" "stackoverflow.com" "npmjs.com" "golang.org")
    
    for domain in "${domains[@]}"; do
        nslookup $domain > /dev/null 2>&1 &
        dig $domain > /dev/null 2>&1 &
    done
    wait
    echo "✅ DNS traffic generated"
}

# Function to generate ping traffic
generate_ping_traffic() {
    echo "🏓 Generating ping traffic..."
    hosts=("8.8.8.8" "1.1.1.1" "google.com")
    
    for host in "${hosts[@]}"; do
        ping -c 5 $host > /dev/null 2>&1 &
    done
    wait
    echo "✅ Ping traffic generated"
}

# Function to test local traffic
generate_local_traffic() {
    echo "🏠 Generating local traffic..."
    # Test API endpoints
    for i in {1..5}; do
        curl -s http://localhost:8080/api/v1/health > /dev/null 2>&1 &
        curl -s http://localhost:8080/api/v1/interfaces > /dev/null 2>&1 &
        curl -s http://localhost:8080/api/v1/stats > /dev/null 2>&1 &
        sleep 0.2
    done
    wait
    echo "✅ Local traffic generated"
}

# Main execution
echo "Select traffic type to generate:"
echo "1) All traffic types"
echo "2) HTTP traffic only"
echo "3) DNS traffic only"
echo "4) Ping traffic only"
echo "5) Local API traffic only"
echo "6) Continuous traffic (press Ctrl+C to stop)"

read -p "Enter your choice (1-6): " choice

case $choice in
    1)
        generate_http_traffic
        generate_dns_traffic
        generate_ping_traffic
        generate_local_traffic
        ;;
    2)
        generate_http_traffic
        ;;
    3)
        generate_dns_traffic
        ;;
    4)
        generate_ping_traffic
        ;;
    5)
        generate_local_traffic
        ;;
    6)
        echo "🔄 Starting continuous traffic generation..."
        echo "Press Ctrl+C to stop"
        while true; do
            generate_http_traffic
            generate_dns_traffic
            generate_local_traffic
            echo "💤 Waiting 10 seconds before next round..."
            sleep 10
        done
        ;;
    *)
        echo "❌ Invalid choice"
        exit 1
        ;;
esac

echo "🎉 Traffic generation completed!"
echo "Check your PCAP Tracking Dashboard at http://localhost:3000"
