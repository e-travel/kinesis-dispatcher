require 'socket'

class UnixDatagramClient

  def initialize
    @socket = Socket.new Socket::PF_UNIX, Socket::SOCK_DGRAM
  end

  def connect
    @socket.connect Socket.sockaddr_un("/tmp/datagram.sock")
  end

  # see http://www.rubydoc.info/stdlib/socket/1.9.3/Socket/Constants
  def send_message(message)
    @socket.send(message, 0)
  end
end

client = UnixDatagramClient.new
client.connect

while message = gets
  client.send_message(message.strip)
end
