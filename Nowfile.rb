
def greet(name = 'World')
    puts "Hello, #{name}!"
end

def nop(apple, banana = 'yes', cucumber:, durian: 'smelly')
    puts apple: apple, banana: banana, cucumber: cucumber, durian: durian
end
