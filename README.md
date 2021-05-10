# Okay

A command runner.

## Usage

Given a Okayfile in any supported language:
```ruby
# Okayfile.rb

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
build                                                     Makefile    make
example <apple> <banana=yellow> --cherry --durian=smelly  Okayfile.rb  ruby
generate                                                  Makefile    make
get <url>                                                 Okayfile.go  go
greet <name=World>                                        Okayfile.rb  ruby
list                                                      Okayfile.go  go
types                                                     Okayfile.go  go
```

## Current Supported Languages/Tools
- Go
- Make
- Ruby
- Yarn

## To do

- [ ] Improve tool interface
- [ ] Tests
- [ ] Scour error paths
- [ ] Add more tools
  - [x] Go
  - [ ] sh
  - [ ] Python
  - [ ] Node
  - [ ] Rake
- [ ] Set up goreleaser
- [ ] Tool inits
- [ ] Help
- [ ] README
- [ ] Cache in temp file?

## Attributions

- [github.com/radovskyb/watcher](https://github.com/radovskyb/watcher)
