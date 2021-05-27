# 👌 `ok`

A task runner.

## Installation

### Brew

```shell
$ brew tap broothie/ok && brew install ok
```

### Releases

Releases can be found on the [releases page](https://github.com/broothie/okay/releases).

### Via Go

```shell
$ go install github.com/broothie/ok/ok.go@latest
```

## Usage

Given an `Okfile` written in any [supported language/tool](#currently-supported-languagestools):
```ruby
# Okfile.rb

def example(apple, banana = 'yellow', cherry:, durian: 'smelly')
    puts "#{apple} apple, #{banana} banana, #{cherry} cherry, #{durian} durian"
end
```

You can use `ok` to call methods directly from the command line:
```shell
$ ok example 'granny smith' --durian stinky -c maraschino
granny smith apple, yellow banana, maraschino cherry, stinky durian
```

You can also run `ok` without a task name to list available tasks:
```shell
$ ok
build                                                     Makefile   make
example <apple> <banana=yellow> --cherry --durian=smelly  Okfile.rb  ruby
generate                                                  Makefile   make
get <url>                                                 Okfile.go  go
greet <name=World>                                        Okfile.rb  ruby
list                                                      Okfile.go  go
types                                                     Okfile.go  go
```

## Currently Supported Languages/Tools
- Go
- Make
- Ruby
- Rake
- Python
- Node
- zsh
- bash
- docker-compose
- Yarn
