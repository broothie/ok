# `now`

A command runner.

## Usage

Given a Nowfile in any supported language:
```ruby
# Nowfile.rb

def example(apple, banana = 'yellow', cherry:, durian: 'smelly')
    puts "#{apple} apple, #{banana} banana, #{cherry} cherry, #{durian} durian"
end
```

You can use `now` to call methods directly from the command line:
```bash
$ now example 'granny smith' --durian stinky -c maraschino
granny smith apple, yellow banana, maraschino cherry, stinky durian
```

You can call `now` without a task name to list available tasks:
```bash
$ now
build                                                         Makefile    make
example <apple> <banana='yellow'> --cherry --durian='smelly'  Nowfile.rb  ruby
get <url>                                                     Nowfile.go  go
greet <name='World'>                                          Nowfile.rb  ruby
types                                                         Nowfile.go  go
```

## To do

- [ ] Improve tool interface
- [ ] Tests
- [ ] Scour error paths
- [ ] Add more tools
  - [ ] Go
  - [ ] sh
  - [ ] Python
  - [ ] Node
- [ ] Set up goreleaser
- [ ] Tool inits
- [ ] Help
- [ ] README
- [ ] Cache in temp file?

## Attributions

- [watcher](https://github.com/radovskyb/watcher)
