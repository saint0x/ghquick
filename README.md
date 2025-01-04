# ghquick 🚀

A lightning-fast CLI tool for GitHub operations with AI-powered automation.

## Features

- 🤖 AI-powered commit message generation
- ⚡ Fast parallel operations
- 🔄 Automatic repository creation and setup
- 🔒 Secure authentication handling
- 🔍 Smart diff detection
- 🔁 Automatic retry mechanism
- 📝 Beautiful logging with progress indicators

## Installation

```bash
# Clone the repository
git clone https://github.com/saint0x/ghquick.git

# Build and install
cd ghquick
./dev.sh
```

## Configuration

Set up the following environment variables in your shell configuration (e.g., ~/.zshrc):

```bash
export GITHUB_TOKEN="your_github_token"
export GITHUB_USERNAME="your_github_username"
export OPENAI_API_KEY="your_openai_api_key"
```

## Usage

### Quick Push with AI-Generated Commit Message

```bash
ghquick push start
```

### Push with Custom Commit Message

```bash
ghquick push --commitmsg "your commit message"
```

### Push to Specific Repository

```bash
ghquick push --name repo-name --commitmsg "your commit message"
```

### Create Private Repository

```bash
ghquick push --name repo-name --private --commitmsg "initial commit"
```

### Debug Mode

```bash
ghquick push start --debug
```

### Custom Timeout

```bash
ghquick push start --timeout 5m
```

## Features in Detail

### AI-Powered Commit Messages
- Uses GPT-4 to analyze your changes
- Generates conventional commit messages
- Understands code context

### Smart Git Operations
- Automatic repository initialization
- Secure credential handling
- Detects and cleans stale locks
- Checks for unpushed changes
- Retries on failure

### Performance Features
- Parallel operations where possible
- Efficient caching
- Smart timeouts
- Optimized for speed

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

MIT License - feel free to use and modify! 