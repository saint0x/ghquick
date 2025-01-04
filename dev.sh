#!/bin/zsh

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo "${GREEN}Installing dependencies...${NC}"
go mod download
go mod tidy

echo "${GREEN}Building ghquick...${NC}"
mkdir -p bin
go build -o bin/ghquick

if [ $? -eq 0 ]; then
    echo "${GREEN}Build successful!${NC}"
    echo "${GREEN}Installing to /usr/local/bin...${NC}"
    cp bin/ghquick /usr/local/bin/
    echo "${GREEN}Done! You can now use 'ghquick' command${NC}"
else
    echo "${RED}Build failed${NC}"
    exit 1
fi 