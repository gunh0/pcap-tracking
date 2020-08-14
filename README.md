# Reference for Network Packet Detection and Analysis Environment

- [How the NSA Can Crack Tor's Anonymity](#How-the-NSA-Can-Crack-Tor's-Anonymity)
- [Demonstration Code Fragments about .pcap](#Demonstration-Code-Fragments-about-.pcap)
  - [dpkt](#dpkt)
  - [pyshark](#pyshark)
- [Compression Method](#Compression-Method)
  - [Zstandard](#Zstandard)
- [How To Transfer Files With Rsync Over SSH](#How-To-Transfer-Files-With-Rsync-Over-SSH)
  - Differences between rsync on remote and rsync local on mounted sshfs?
  - [rsync from Linux to Windows over SSH](#rsync-from-Linux-to-Windows-over-SSH)

<br/>

-----

<br/>

## How the NSA Can Crack Tor's Anonymity

Tor has always kept web traffic anonymous by delaying or otherwise altering the packets of data that are sent through servers (that's why Tor tends to run slower than your standard browser), making it look like the traffic is coming from a place that it's not actually (the IP address that the server "sees" is called an "exit node"). If the end server of the site you visited can also detect the origin point of where your traffic enter Tor (the "entry node" or "entry relay), then anonymity is lost.

![image](https://user-images.githubusercontent.com/41619898/85376469-b6715000-b572-11ea-8719-ef729b687dfa.png)

------

<br/>

## Demonstration Code Fragments about .pcap

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

<br/>

-----

<br/>

## Compression Method

윈도우즈에서의 압축은 zip 등의 방식으로 파일이나 폴더들을 묶음과 동시에 압축하는 것을 의미하지만,

리눅스에서는 파일이나 폴더들을 묶는 것(archive)과 실제로 압축(compress)하는 기능으로 세분화 되어 있다.

<br/>

리눅스에서 여러 파일을 한 파일로 묶는 것을 아카이브라하며 확장자는 `.tar`(tape archives)이다.

그리고 보통 생성된 tar 파일을 다시 gzip을 사용하여 압축해서 .tar.gz의 확장자를 가지는 압축 아카이브파일을 많이 사용한다.

![image-20200707141424833](README.assets/image-20200707141424833.png)

![image-20200707141723851](README.assets/image-20200707141723851.png)

<br/>

### Zstandard

**Zstandard** (also known as **zstd**) is a free open source, fast real-time data compression program with better compression ratios, developed by **Facebook**. It is a lossless compression algorithm written in **C** (there is a re-implementation in **Java**) – its thus a native Linux program.

The following results are obtained by doing several fast compression algorithms tests on a server running Linux Debian using [lzbench](https://github.com/inikep/lzbench), an open-source in-memory benchmark tool.

![img](https://www.tecmint.com/wp-content/uploads/2018/06/Zstandard-Compression-Testing.png)

<br/>

##### How to Install Zstandard Compression Tool in Linux

```bash
$ sudo apt update && sudo apt install build-essential		#Ubuntu/Debian
# yum group install "Development Tools" 					#CentOS/REHL
# dnf groupinstall "C Development Tools and Libraries"		#Fedora 22+
```

<br/>

**Reference:** https://github.com/facebook/zstd

<br/>

----

<br/>

## How To Transfer Files With Rsync Over SSH

**Reference:** https://phoenixnap.com/kb/how-to-rsync-over-ssh	| Posted January 31, 2020

<br/>

##### Verify Rsync Installation

Most recent Linux distributions have **rsync** by default. To verify your system has rsync installed, run the installation command.

For **Debian and Ubuntu** machines, use **`apt-get`**:

```output
sudo apt-get install rsync
```

If rsync is already installed, the output shows the current version of the tool.

On **RPM-based** machines, such as CentOS, use:

```output
sudo yum install rsync
```

<br/>

##### Transfer Files with Rsync over SSH

Before you can start transferring files and directories with rsync over SSH, make sure you can [use SSH to connect to a remote server](https://phoenixnap.com/kb/ssh-to-connect-to-remote-server-linux-or-windows). Once verified, you can begin backing up your data. Ensure your destination system has sufficient storage space.

The syntax for copying files to a remote server over SSH with the **`rsync`** command is:

```output
rsync OPTION SourceDirectory_or_filePath user@serverIP_or_name:Target
```

<br/>

##### Check Rsync File Transfer Progress

To check the status of rsync transfers, use the **`-P`** option. This option displays the transfer times, as well as the names of the files and directories that are synced.

If there is an issue with your connection and the sync is interrupted, **`-P`** resumes your transfers.

Run the command in this format to sync recursively and check the status of the transfer:

```output
rsync -aP ~/SourceDirectory/ username@192.168.56.100:~/Destination
```

<br/>

### Differences between rsync on remote and rsync local on mounted sshfs?

> SSHFS is convenient, but it doesn't mesh well with rsync or more generally with synchronization tools.
>
> The biggest problem is that SSHFS largely kills rsync's performance optimizations. In particular, for medium to large files, when rsync sees that a file has been modified, it calculates checksums on parts of the file on each side in order to transfer only the parts that have been modified. This is an optimization only if the network bandwidth is significantly smaller than the disk bandwidth, which is usually the case. But with SSHFS, the “disk” bandwidth is in fact the network bandwidth, so rsync would have to read the whole file in order to determine what to copy. In fact, with a local copy (which it is, as far as rsync is concerned, even if one of the sides is on SSHFS), rsync just copies the whole file.
>
> SSHFS is also detrimental to performance if there are many small files. Rsync needs to check at least the metadata of every file to determine whether it's been modified. With SSHFS, this requires a network round trip for each file. With rsync over SSH, the two sides can work in parallel and transfer information in bulk, which is a lot faster.
>
> In terms of access restrictions, SSHFS requires SFTP access, whereas rsync requires the ability to run code (specifically, the rsync program) via a shell. If the user doesn't have a shell account, It's possible and common to provide an account with a special shell that only allows running a few programs including `sftp-server` and `rsync`.
>
> Reference: https://unix.stackexchange.com/questions/283838/differences-between-rsync-on-remote-and-rsync-local-on-mounted-sshfs

<br/>

### rsync from Linux to Windows over SSH

Cygwin is a POSIX-compatible programming and runtime environment that runs natively on Microsoft Windows. Under Cygwin, source code designed for Unix-like operating systems may be compiled and run natively with minimal modification.

> I ended up installing [cygwin](https://cygwin.com/setup-x86_64.exe), and made sure to also install the `rsync` package. Then, I modified my `PATH` environment variable to include cygwin's `bin` directory. I was then able to call `rsync` from powershell, to confirm the installation worked.
>
> After that, I was able to successfully use `rsync` on my Linux machine to transfer files to the Windows machine.

<br/>

------

<br/>

