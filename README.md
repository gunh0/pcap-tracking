### How the NSA Can Crack Tor's Anonymity

Tor has always kept web traffic anonymous by delaying or otherwise altering the packets of data that are sent through servers (that's why Tor tends to run slower than your standard browser), making it look like the traffic is coming from a place that it's not actually (the IP address that the server "sees" is called an "exit node"). If the end server of the site you visited can also detect the origin point of where your traffic enter Tor (the "entry node" or "entry relay), then anonymity is lost.

![image](https://user-images.githubusercontent.com/41619898/85376469-b6715000-b572-11ea-8719-ef729b687dfa.png)

------

<br/>

## Demonstration Code Fragments for Solutions

```powershell
> pip install virtualenv
> virtualenv venv
> .\venv\Scripts\activate
(venv) > python -m pip install -r requirements.txt
```

<br/>

### dpkt

A fast, simple packet creation/parsing python module with definitions for the basic TCP/IP protocols.

- https://dpkt.readthedocs.io/en/latest/print_packets.html

<br/>

### pyshark

Python wrapper for tshark, allowing python packet parsing using wireshark dissectors.

There are quite a few python packet parsing modules, this one is different because it doesn't actually parse any packets, it simply uses tshark's (wireshark command-line utility) ability to export XMLs to use its parsing.

This package allows parsing from a capture file or a live capture, using all wireshark dissectors you have installed. Tested on windows/linux.

>Pyshark is a packet parsing only package, being a wrapper of tshark.

<br/>

##### transport_layer

```python
def network_conversation(packet):
    try:
        protocol = packet.transport_layer
        source_address = packet.ip.src
        source_port = packet[packet.transport_layer].srcport
        destination_address = packet.ip.dst
        destination_port = packet[packet.transport_layer].dstport
        return (f'{protocol} {source_address}:{source_port} ==> {destination_address}:{destination_port}')
    except AttributeError as e:
        pass
```

<br/>

##### highest_layer

```python
def network_conversation(packet):
    try:
        protocol = packet.highest_layer
        source_address = packet.ip.src
        source_port = packet[packet.highest_layer].srcport
        destination_address = packet.ip.dst
        destination_port = packet[packet.highest_layer].dstport
        return (f'{protocol} {source_address}:{source_port} ==> {destination_address}:{destination_port}')
    except AttributeError as e:
        pass
```

