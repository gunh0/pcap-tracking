import pyshark
import os
import sys
from scapy.all import *

capture = pyshark.FileCapture('res/Tor_Starting.pcap')

filename = 'New_Tor_Starting.pcap'
new_cap = PcapWriter(filename, append=True)
before_num = 1
after_num = 1
for packet in capture:
    before_num += 1
    try:
        protocol = packet.transport_layer
        print(before_num, after_num, packet.transport_layer, packet.ip.src, packet.ip.dst)
        after_num += 1
        #new_cap.write(packet)
    except AttributeError as e:
        pass