# Web3Safe

Web3Safe is a set of command-line tools designed to protect your development
environment by analyzing shell environment variables and .env (dotenv) files
for any sensitive information, such as `PRIVATE_KEY`, `MNEMONIC`, and many
other variables that can be stolen by malware or degens.

## Personal Story

Web3Safe was created from a personal experience that showed how important it
is to keep our data safe while working on projects.

Long story short: one day, I got a message on LinkedIn asking for help with a
web3 app. I was excited to help and started working on it right away.

But then something unexpected happened. The project had hidden obfuscated code
that secretly looked through all my files, including sensitive ones like .env
files. Before I knew it, I lost access to my wallet and tokens.

That's why I made Web3Safe. It's a tool that helps developers like us keep our
work safe. With Web3Safe, you can check your computer for any problems with
your environment variables and make sure your projects stay secure.

## Features

- Analyzes shell environment variables for sensitive information.
- Scans .env files for sensitive data such as passwords, API keys, and other confidential information.
- Scans all keys in YAML files for sensitive data for sensitive information.
- Provides customizable and extendable rules.
- Supports exclusion of certain files from the analysis.

## Getting Started

### Installation

Web3Safe is a command-line tool written in Go. To install it, follow these steps:

1. Clone the repository:
   ```
   git clone https://github.com/gruz0/web3safe.git
   ```
2. Build the apps:
   ```
   cd web3safe
   make build
   ```

3. App will be placed inside `bin` directory:
   ```
   web3safe
   ```

### Docker

TBD

## Usage

### Create a new configuration file

```sh
web3safe config -create [-config "/path/to/config.yml"] [-force]
```

### Print the default config (or a given config) to your terminal

```sh
web3safe config -print [-config "/path/to/config.yml"]
```

### Analyze shell ENV variables

This tool scans the current user's shell environment variables and display any
sensitive information found.

```sh
web3safe shellenv [-config "/path/to/config.yml"]
```

Example:

```sh
$ MNEMONIC=test web3safe shellenv

Shell ENV has a sensitive variable: MNEMONIC
```

### Analyze dotenv (.env) files

```sh
web3safe dotenv [-config "/path/to/config.yml"]
```

You can also customize the analysis by providing additional flags:

- `-dir`: Path to the directory to scan
- `-recursive`: If set, the directory will be scanned recursively
- `-file`: Path to the file to scan

Example:

```sh
$ web3safe dotenv -dir . -recursive

samples/.env:5: found sensitive variable MNEMONIC_WORDS
samples/.env:7: found sensitive variable private_key
samples/.env.export:1: found sensitive variable PRIVATE_KEY
samples/.env.export:2: found sensitive variable BINANCE_ACCOUNT_PRIVATE_KEY
```

### Analyze YAML files

```sh
web3safe yaml [-config "/path/to/config.yml"]
```

You can also customize the analysis by providing additional flags:

- `-dir`: Path to the directory to scan
- `-recursive`: If set, the directory will be scanned recursively
- `-file`: Path to the file to scan

Example:

```sh
$ web3safe yaml -dir . -recursive

samples/config.yml: found sensitive key "PASSWORD" in .nested.inside.PASSWORD
samples/config.yml: found sensitive key "MNEMONIC" in .nested.inside.MNEMONIC
samples/playbook.yml: found sensitive key "password" in [0].password
samples/playbook.yml: found sensitive key "mnemonic" in [0].env.mnemonic
```

## Contributing

Contributions to Web3Safe are welcome! If you encounter any bugs, issues, or
have suggestions for improvement, please open an issue on GitHub or submit a
pull request with your changes.

## License

Web3Safe is licensed under the MIT License. Feel free to use, modify,
and distribute the code for both commercial and non-commercial purposes.
