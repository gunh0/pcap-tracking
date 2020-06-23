import pyshark

def network_trans_conversation(packet):
    try:
        protocol = packet.transport_layer
        source_address = packet.ip.src
        source_port = packet[packet.transport_layer].srcport
        destination_address = packet.ip.dst
        destination_port = packet[packet.transport_layer].dstport
        return (f'{protocol} {source_address}:{source_port} ==> {destination_address}:{destination_port}')
    except AttributeError as e:
        pass


def network_high_conversation(packet):
    try:
        protocol = packet.highest_layer
        source_address = packet.ip.src
        source_port = packet[packet.highest_layer].srcport
        destination_address = packet.ip.dst
        destination_port = packet[packet.highest_layer].dstport
        return (f'{protocol} {source_address}:{source_port} ==> {destination_address}:{destination_port}')
    except AttributeError as e:
        pass


capture = pyshark.FileCapture('res/Tor_Starting.pcap')
conversations = []
for packet in capture:
    results = network_trans_conversation(packet)
    if results != None:
        conversations.append(results)

num = 1
for item in conversations:
    print(num, item)
    num += 1

'''
# this sorts the conversations by protocol 
for item in sorted(conversations):
  print (item)
'''
