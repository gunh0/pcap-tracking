import socket
import threading
import os

def handle_client_connection(client_socket):
    request = client_socket.recv(1024)
    print('Received {}'.format(request))
    ack = os.system(request)
    print(request)
    if "ERROR" in str(ack) and "Connection timed out" in str(ack):
        print("CLOSING...Because of CONNECTION ERROR" + "\n")
        client_socket.send(bytes(ack))
        client_socket.close()
        return 1
    else:
        print("ack :: "+str(ack)+"\n")
        client_socket.send(bytes(ack))
        client_socket.close()
        return 1


if __name__ == "__main__":
    bind_ip = '0.0.0.0'
    bind_port = 9999

    server = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    server.bind((bind_ip, bind_port))
    server.listen(5)

    print('Listening on {}:{}'.format(bind_ip, bind_port))

    while True:
        client_sock, address = server.accept()
        print('Accepted connection from {}:{}'.format(address[0], address[1]))
        client_handler = threading.Thread(
            target=handle_client_connection,
            args=(client_sock,)
        )
        client_handler.start()
