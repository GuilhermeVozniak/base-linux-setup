#!/bin/bash

# Base Linux Setup - Example Usage Script
# This script demonstrates how to use the base-linux-setup tool

set -e  # Exit on any error

echo "🚀 Base Linux Setup - Example Usage"
echo "=================================="
echo

# Check if the tool is built
if [ ! -f "build/base-linux-setup" ]; then
    echo "📦 Building the application..."
    make build
    echo
fi

echo "📋 Available Commands:"
echo "1. Help - Show all available commands"
echo "2. List Presets - Show all available presets"
echo "3. Detect Environment - Detect current system"
echo "4. Run Setup - Interactive setup process"
echo

read -p "Enter your choice (1-4): " choice

case $choice in
    1)
        echo "🔍 Showing help..."
        ./build/base-linux-setup --help
        ;;
    2)
        echo "📋 Listing all presets..."
        ./build/base-linux-setup list-presets
        ;;
    3)
        echo "🔍 Detecting environment..."
        echo "Note: This requires neofetch to be installed"
        ./build/base-linux-setup detect
        ;;
    4)
        echo "🚀 Starting interactive setup..."
        echo "Note: This will attempt to detect your environment and run setup"
        echo "Warning: This may make changes to your system!"
        read -p "Are you sure you want to continue? (y/N): " confirm
        if [[ $confirm == [yY] ]]; then
            ./build/base-linux-setup
        else
            echo "Setup cancelled."
        fi
        ;;
    *)
        echo "❌ Invalid choice. Please run the script again."
        exit 1
        ;;
esac

echo
echo "✅ Example completed!"
echo
echo "💡 Quick Start Guide:"
echo "• Install neofetch first: sudo apt-get install neofetch"
echo "• Run: ./build/base-linux-setup"
echo "• Follow the interactive prompts"
echo "• Customize tasks as needed"
echo "• Review and confirm before execution"
echo
echo "🔧 Development Commands:"
echo "• make build     - Build the application"
echo "• make install   - Install to /usr/local/bin"
echo "• make test      - Run tests"
echo "• make clean     - Clean build artifacts"
echo "• make help      - Show all make targets" 