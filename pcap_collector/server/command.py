import socket
import sys
import os
import queue
import datetime


def command(ip, port, command):
    # create an ipv4 (AF_INET) socket object using the tcp protocol (SOCK_STREAM)
    client = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    # connect the client
    # client.connect((target, port))
    try:
        client.connect((ip, int(port)))
    except socket.error:
        print(ip, ":", port, " - CONNECTION ERROR")
        client.close()
        return
    # send some data (in this case a HTTP GET request)
    if "windows-tor" in command:
        client.send(command.encode())
        return
    if "change_name" in command:
        temp = ip[-3:]
        command = command.replace('change_name', 'Tor-%s-ymd' % temp)
        date = str(datetime.datetime.now())
        date = date[2:19]
        date = date.replace('-', '').replace(' ', '').replace(':', '')
        # hostname = socket.gethostname()
        # name = '%s-%s' % (hostname, date)
        command = command.replace('ymd', '%s' % date)
        client.send(command.encode())
        return

    if "Bye" in command:
        print("Target ip = %s \nRunning command = %s\n" % (ip, command))
        temp = "close"
        client.send(temp.encode())
        temp = client.recv(4096)
        temp = temp.decode()
        client.close()
        return temp
    else:
        client.send(command.encode())
    # receive the response data (4096 is recommended buffer size)
    # result = (int.from_bytes(client.recv(4096),byteorder='little'))
    # if "watch.py" in command:
    #     temp = "sudo python3 /home/tor/Tor_CIFS/Browser_Crawler/watch.py"
    #     print("Target ip = %s \nRunning command = %s\n"%(ip, temp))
    #     return

    print("Target ip = %s \n Running command = %s\n" % (ip, command))
    # if "capture.sh" in command:
    #     print("DODO")
    #     return
    # else:
    #     print("RUN")
    temp = client.recv(4096)
    temp = temp.decode()
    return temp


def search(dirname):
    filenames = os.listdir(dirname)
    files = []
    for filename in filenames:
        full_filename = os.path.join(dirname, filename)
        ext = os.path.splitext(full_filename)[-1]
        if ext == '.pcap':
            files.append(full_filename)
    return files


def cmd_cmd_list(send_com):
    cmd_list = queue.Queue()
    cmd_list.put(send_com)
    return cmd_list


def capture(ip, url, path):
    date = str(datetime.datetime.now())
    date = date[2:19]
    date = date.replace('-', '').replace(' ', '').replace(':', '')
    # hostname = socket.gethostname()
    # name = '%s-%s' % (hostname, date)
    interface = 'eth0'
    url = url.replace("http://", '')
    make = 'mkdir %s%s' % (path, url)
    if ip == 0:
        com = '''tshark -i %s -w %s%s/change_name.pcap ''' % (
            interface, path, url)
    else:
        com = '''tshark -i %s -f "host %s" -w %s%s/change_name.pcap ''' % (
            interface, ip, path, url)
    return make, com


def make_moz_cmd_list(checker, ip, url, level, path):
    cmd_list = queue.Queue()
    # url_replace = url.replace("http://", '').replace("https://",'')

    make, tshark_on = capture(ip, url, path)
    cmd_list.put('mkdir %s' % path)
    cmd_list.put(make)
    cmd_list.put(tshark_on)
    # tshark_on = "/home/tor/Tor_CIFS/Source/scripts/capture.sh " + url_replace
    # cmd_list.put(tshark_on)

    if checker == 1:
        collect_cmd = "wget -P /home/tor/trash -E -k -p -e robots=off --tries=3 %s" % url
        cmd_list.put(collect_cmd)
    elif checker == 2:
        # url_replace = url[37:] Using on Copy
        make, tshark_on = capture(ip, url, path)
        cmd_list.put('mkdir %s' % path)
        cmd_list.put(make)
        cmd_list.put(tshark_on)
        collect_cmd = "torify wget -P /home/tor/trash -E -k -p -e robots=off --tries=3 %s" % url
        cmd_list.put(collect_cmd)
    # torify = "torify wget -P /home/tor/trash -E -k -p -e robots=off "  + url
    # torify = "wget -P /home/tor/trash -E -k -p -e robots=off --tries=3 " + url
    # cmd_list.put(torify)

    rm = "rm -rf /home/tor/trash/*"
    cmd_list.put(rm)

    tshark_off = "killall -9 dumpcap"
    cmd_list.put(tshark_off)
    # print(cmd_list)
    # pcap2csv = "./pcap2csv.py " + url + " " + path
    # cmd_list.put(pcap2csv)

    return cmd_list


def make_tbb_cmd_list(address_set):
    cmd_list = queue.Queue()
    collect_cmd = "sudo python3 /home/jjangga94temp/Tor_CIFS/Browser_Crawler/TBB.py -u %s" % address_set
    cmd_list.put(collect_cmd)
    print(collect_cmd)
    return cmd_list


def watch_cmd_list(address_set):
    cmd_list = queue.Queue()
    watch_cmd = "sudo python3 /home/jjangga94temp/Tor_CIFS/Browser_Crawler/watch.py -u %s" % address_set
    cmd_list.put(watch_cmd)
    print(cmd_list)
    return cmd_list


def make_copy_cmd_list():
    copy_cmd_list = queue.Queue()
    copy_pcap = "cp -r /home/tor/tor/0625/* /home/jjangga94temp/Tor_CIFS/Traffic/0717"
    # copy_pcap = "cp -r /home/tor/tor/100_Hidden_TBB_0310_non_fix/* /home/jjangga94temp/Tor_CIFS/Traffic/2020_Traffic"
    # copy_pcap = "sudo scp /home/tor/tor/100_Hidden_TBB_0310_non_fix/* 127.0.0.1:/home/jjangga94temp/Tor_CIFS/Traffic/fastTest"
    # copy_pcap = "cp -r /home/tor/tor/100_Hidden_TBB_0310_non_fix/* /home/jjangga94temp/Tor_CIFS/Traffic/0310_non_fix"
    # copy_pcap = "cd /home/tor/tor/100_Hidden_TBB_0310_non_fix/; tar cz * | ssh 121.67.187.138 tar xz -C /home/jjangga94temp/Tor_CIFS/Traffic/fastTest"
    # copy_pcap = "cp -r /home/tor/tor/100_Hidden_TBB/* /home/tor/Tor_CIFS/Json"
    copy_cmd_list.put(copy_pcap)
    print(copy_cmd_list)
    # pcap_dir_list = os.listdir('/RESUL_TEST/results')
    # for this in pcap_dir_list:
    #     path = os.path.abspath(this)
    #     path_basename = os.path.basename(path)
    #     copy_pcap = "cp -r %s /home/tor/Tor_CIFS/Traffic/%s"%(path,path_basename)
    #     copy_cmd_list.put(copy_pcap)
    return copy_cmd_list


def exit_cmd_list():
    exit_cmd_list = queue.Queue()
    exit = "Bye"
    exit_cmd_list.put(exit)
    return exit_cmd_list
