import socket

s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s.connect(("127.0.0.1", 8080))

request1 = (
    "GET /first HTTP/1.1\r\n"
    "Host: localhost\r\n"
    "\r\n"
)

request2 = (
    "GET /second HTTP/1.1\r\n"
    "Host: localhost\r\n"
    "\r\n"
)

# Send both requests back-to-back on the same TCP connection
s.sendall((request1 + request2).encode())

# Read whatever response your server sends
while True:
    data = s.recv(4096)
    if not data:
        break

    print(data.decode(), end="")

s.close()