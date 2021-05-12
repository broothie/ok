require 'dotenv/load'
require 'pry'

def publish(message)
    bump
    puts `git add -A`
    puts `git commit -m "{message}"`
    puts `git push`
end

def bump
    puts `bump ok/version.go`
end

def fix_imports
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

def example(apple, banana = 'yellow', cherry:, durian: 'smelly')
    puts "#{apple} apple, #{banana} banana, #{cherry} cherry, #{durian} durian"
end

def repeat(message, n = 3)
    n.times { puts message }
end

def greet(name = 'andrew')
    binding.pry
end
