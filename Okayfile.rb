require 'pry'

def example(apple, banana = 'yellow', cherry:, durian: 'smelly')
    puts "#{apple} apple, #{banana} banana, #{cherry} cherry, #{durian} durian"
end

def repeat(message, n = 3)
    n.times { puts message }
end

def fix_imports
    filenames = Dir['**/*.go']
    filenames
end

def repl
    binding.pry
end
