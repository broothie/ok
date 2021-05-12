require 'dotenv/load'
require 'pry'

def release(message, dry: true)
    bump
    version = get_version
    tag version, message
    push_tag version

    if dry
        puts `goreleaser --snapshot --skip-publish --rm-dist`
    else
        puts `goreleaser release`
    end
end

def publish(message)
    bump
    version = get_version
    tag version, message
    push_tag version
end

def tag(version, message)
    puts `git tag -a #{version} -m "#{message}"`
end

def push_tag(version)
    puts `git push origin #{version}`
end

def bump
    puts `bump ok/version.go`
end

def get_version
    p File.read('ok/version.go').match(/v\d+\.\d+\.\d+/)[0]
end

def example(apple, banana = 'yellow', cherry:, durian: 'smelly')
    puts "#{apple} apple, #{banana} banana, #{cherry} cherry, #{durian} durian"
end

def repeat(message, n = 3)
    n.times { puts message }
end

def greet(name = 'andrew')
    binding.pry
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
