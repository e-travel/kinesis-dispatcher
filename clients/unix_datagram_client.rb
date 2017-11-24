require 'socket'

socket_address = ARGV.last || '/tmp/my.sock'
puts "set address=#{socket_address}"

class UnixDatagramClient
  def initialize(address)
    @socket = Socket.new Socket::PF_UNIX, Socket::SOCK_DGRAM
    @address = address.strip
  end

  def connect
    @socket.connect Socket.sockaddr_un(@address)
  rescue
    puts "Failed to connect to #{@address}. Aborting."
    exit 1
  end

  # see http://www.rubydoc.info/stdlib/socket/1.9.3/Socket/Constants
  def send_message(message)
    @socket.send(message, 0)
  end
end

client = UnixDatagramClient.new socket_address
client.connect

puts "Enter messages:"
while message = gets
  begin
    client.send_message(message.strip)
  rescue
    puts "Send Failed. Reconnecting..."
    client.connect
  end
end
