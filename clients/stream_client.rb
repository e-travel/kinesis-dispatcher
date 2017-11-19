require 'socket'

type = ARGV.last
if type.nil?
  puts("please specify unix or tcp")
  exit(1)
end

def send_message(type, content)
  socket =
    case type
    when 'unix'
      UNIXSocket.new '/tmp/message-dispatcher.sock'
    when 'tcp'
      TCPSocket.new 'localhost', 8888
    else
      raise "type %s not supported"
    end
  socket.write(content)
  socket.close
end

def messages(type)
  m1 = "Hello there\nHello"
  m2 = "Goodbye man\nThat was good"
  send_message(type, m1)
  send_message(type, m2)
end


messages(type)
