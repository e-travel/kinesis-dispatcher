require 'socket'

type = ARGV.last
if type.nil?
  puts("please specify unix or tcp")
  exit(1)
end

puts "Each new line is a new message. C-d to finish."

socket =
  case type
  when 'unix'
    UNIXSocket.new '/tmp/kinesis-dispatcher.sock'
  when 'tcp'
    TCPSocket.new 'localhost', 8888
  else
    raise "type %s not supported"
  end

while (msg = STDIN.gets)
  socket.write(msg)
end

socket.close
