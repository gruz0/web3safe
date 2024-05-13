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

3. All apps will be placed inside `bin` directory:
   ```
   dotenvanalyzer
   envanalyzer
   yamlanalyzer
   web3safe
   ```

Web3Safe includes three tools: Web3Safe itself, Shell ENV Analyzer and Dotenv
Analyzer.

### Web3Safe

This tool is designed for creating a configuration file for the other apps.

- `-help`: Show all the available commands.
- `-generateConfig`: Generate a new configuration file.

### EnvAnalyzer

This tool scans the current user's environment variables and display any
sensitive information found.

You can also customize the analysis by providing additional flags:

- `-help`: Show all the available commands.
- `-config`: Specify a custom configuration file for rule customization.

### DotEnvAnalyzer

By default, this tool scans .env files starting from a given directory
recursively and display any sensitive information found inside `.env` files.

You can also customize the analysis by providing additional flags:

- `-help`: Show all the available commands.
- `-config`: Specify a custom configuration file for rule customization.
- `-path`: Path to start scan from (default: current dir).

### YamlAnalyzer

By default, this tool scans YAML files (`yml` and `yaml`) starting from a given
directory recursively and display any sensitive information found inside files.

You can also customize the analysis by providing additional flags:

- `-help`: Show all the available commands.
- `-config`: Specify a custom configuration file for rule customization.
- `-path`: Path to start scan from (default: current dir).

## Contributing

Contributions to Web3Safe are welcome! If you encounter any bugs, issues, or
have suggestions for improvement, please open an issue on GitHub or submit a
pull request with your changes.

## License

Web3Safe is licensed under the MIT License. Feel free to use, modify,
and distribute the code for both commercial and non-commercial purposes.
