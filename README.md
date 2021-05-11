# 👌 `ok`

A task runner.

## Usage

Given a Okfile in any supported language:
```ruby
# Okfile.rb

def example(apple, banana = 'yellow', cherry:, durian: 'smelly')
    puts "#{apple} apple, #{banana} banana, #{cherry} cherry, #{durian} durian"
end
```

You can use `ok` to call methods directly from the command line:
```bash
$ ok example 'granny smith' --durian stinky -c maraschino
granny smith apple, yellow banana, maraschino cherry, stinky durian
```

You can call `ok` without a task name to list available tasks:
```bash
$ ok
build                                                     Makefile     make
example <apple> <banana=yellow> --cherry --durian=smelly  Okfile.rb  ruby
generate                                                  Makefile     make
get <url>                                                 Okfile.go  go
greet <name=World>                                        Okfile.rb  ruby
list                                                      Okfile.go  go
types                                                     Okfile.go  go
```

## Current Supported Languages/Tools
- Go
- Make
- Ruby
- Node
- docker-compose
- Yarn

## To do

- [ ] Improve tool interface
- [ ] Tests
- [ ] Scour error paths
- [ ] Support .rc or something
- [ ] Task inspect
- [ ] Param validator (validates tool param output)
- [ ] Add more tools
  - [x] Go
  - [ ] sh
  - [ ] Python
  - [x] Node
  - [ ] Rake
  - [x] docker-compose
  - [ ] npm
- [ ] Set up goreleaser
- [x] Tool inits
- [x] Help
- [ ] README
- [ ] Cache in temp file?
- [ ] Tab completion
