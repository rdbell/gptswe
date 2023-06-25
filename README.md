# GPT-SWE

GPT-SWE is a command line tool that assists you in managing your software engineering tasks, leveraging the power of OpenAI's ChatGPT models.

The 'apply changes' feature can create/update/delete files on your system. Use at your own risk. At the very least you should `git commit` anything you want to save before using this feature.

Feel free to submit PRs. Any improvements are appreciated.

## Features

- Interact with ChatGPT to:
  - Ask questions about your code
  - Implement new features
  - Fix bugs
  - Clean up and refactor your code
  - Write tests
  - Find potential bugs
  - Perform code reviews
  - etc.

- Automatically apply changes directly to your project files and let ChatGPT do all of your work. Dangerous and fun!

## Requirements

- An OpenAI API key with access to ChatGPT models.

## Installation

1. Clone the repository:

```sh
$ git clone https://github.com/user/gptswe.git
```

2. Enter the repository and build the executable:

```sh
$ cd gptswe
$ go build
```

3. Set the OpenAI API key as an environment variable:

```sh
$ export OPENAI_API_KEY="yourapikey"
```

## Usage

### Method 1: Binary usage

```sh
$ ./gptswe [file1] [file2] ...
```

```sh
$ ./gptswe *.go
```

Use CLI flags to skip prompts:

```sh
$ ./gptswe --command {command_number} --details "Additional details" [file1] [file2] ...
```

### Method 2: Using Docker

```sh
$ docker run --rm -it -e OPENAI_API_KEY="yourapikey" -v "$(pwd)":/app/files rdbell/gptswe:latest [file1] [file2] ...
```

```sh
$ docker run --rm -it -e OPENAI_API_KEY="yourapikey" -v "$(pwd)":/app/files rdbell/gptswe:latest *.go
```

Use CLI flags to skip prompts:

```sh
$ docker run --rm -it -e OPENAI_API_KEY="yourapikey" -v "$(pwd)":/app/files rdbell/gptswe:latest --command {command_number} --details "Additional details" [file1] [file2] ...
```

## Examples

Run GPT-SWE on your project files:

```sh
$ ./gptswe main.go utils.go
```

Request GPT-SWE to help you implement a new feature in your project:

```sh
$ ./gptswe --command 2 --details "Add login functionality with username and password validation" main.go auth.go
```
