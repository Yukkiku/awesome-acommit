
# Awesome-Acommit

A tool for generating commit messages using the ChatGPT API.

## Installation

Download the binary from [this repository's releases] (https://github.com/Yukkiku/awesome-acommit/releases/) and include it in your $PATH. Alternatively, you can use the `go install` command:

```bash
go install github.com/Yukkiku/awesome-acommit
```

## Usage

```bash
# inside your repo
git add .
acommit
```

## Customization

Modify `~/.config/awesome-acommit/prompt.txt` to customize the prompt.

### Multilingual Support

If you need to write commit messages in a language other than English, simply add your preferred language after the prompt. For instance, for Japanese: