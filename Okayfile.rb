require 'pry'

def example(apple, banana = 'yellow', cherry:, durian: 'smelly')
    puts "#{apple} apple, #{banana} banana, #{cherry} cherry, #{durian} durian"
end

def repeat(message, n = 3)
    n.times { puts message }
end

def fix_imports
    filenames = Dir['**/*.go']
    filenames.each do |filename|
        source = File.read(filename)
        match = source.match(/import\s+\((.*?)\)/im)
        next unless match

        imports = match[1]
        count = imports.chars.each_cons(2).map(&:join).count("\n\n")
        puts filename if count > 1
    end
end
