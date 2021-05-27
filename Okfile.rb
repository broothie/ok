require 'pry'

# List files with more than one double newline
def ugly_imports
    filenames = Dir['**/*.go']
    filenames.each do |filename|
        source = File.read(filename)
        match = source.match(/import\s+\((.*?)\)/m)
        next unless match

        imports = match[1]
        count = imports.chars.each_cons(2).map(&:join).count("\n\n")
        puts filename if count > 1
    end
end

# Examples

# Repeat message `n` times
def repeat(message, n = 3)
    n.times { puts message }
end

# Greet someone by name
def greet(name = 'World')
    puts "Hello, #{name}"
end
