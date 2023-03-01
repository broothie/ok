# `ok`

A task runner.

## Installation

```
go install github.com/broothie/ok
```

## Usage

### Step 1: Define your tasks.

```ruby
# Okfile.rb

def greet(name, greeting: "Hello")
    puts "#{greeting}, #{name}!"
end
```

```go
// Okfile.go
//go:build okfile

package okfile

import "fmt"

func farewell(name string) {
	fmt.Printf("Goodbye, %s\n", name)
}
```

```makefile
# Makefile

build:
    go build -o server cmd/server.go
```

### Step 2: List your tasks.

```
$ ok
TASK      ARGS                     FILE
build                              Makefile
farewell  <name>                   Okfile.go
greet     <name> --greeting=Hello  Okfile.rb
```

### Step 3: Run your tasks.

```
$ ok greet Andrew --greeting Yo
Yo, Andrew!
```
