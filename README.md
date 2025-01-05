# ghquick 🚀

A lightning-fast CLI tool for GitHub operations with AI-powered automation.

## Features

- 🤖 AI-powered commit messages (1-3 words)
- 🔄 Automatic branch creation and PR workflow
- 📝 Interactive PR merging

## Installation

```bash
# Clone the repository
git clone https://github.com/saint0x/ghquick.git

# Build and install
cd ghquick
./dev.sh
```

## Configuration

Set up your environment variables:

```bash
export GITHUB_TOKEN="your_token"
export GITHUB_USERNAME="your_username"
export OPENAI_API_KEY="your_key"  # Required for AI features
```

## Usage

### Push Changes

```bash
# Push with AI-generated commit message (default)
ghquick push

# Push with manual commit message
ghquick push --commitmsg "fix bug"
```

### Create Pull Request
Creates a new branch automatically and opens a PR:

```bash
# Create PR with AI-generated title (default)
ghquick pr create

# Create PR with custom title
ghquick pr create --title "fix bug"

# Create PR against different base branch
ghquick pr create --base develop
```

### Merge Pull Request

```bash
# Interactive PR selection (default)
ghquick pr merge

# Merge specific PR
ghquick pr merge --number 123
```

### Shell Completion

For enhanced CLI experience, `ghquick` supports shell completion in bash/zsh:

```bash
ghquick completion [bash|zsh] > /path/to/completion/script
```

## License

MIT License - feel free to use and modify! 