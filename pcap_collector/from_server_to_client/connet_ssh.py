import os, sys
import paramiko
import csv

gcp_ip_list=[]
user_info=[]
command_list=[]

try:
    with open('private/ip_list.csv') as csv_file:
        csv_reader = csv.reader(csv_file)
        for row in csv_reader:
            gcp_ip_list.append(row[1])
    with open('private/user_info.csv') as csv_file:
        csv_reader = csv.reader(csv_file)
        for row in csv_reader:
            user_info.append(row[0])    # username
            user_info.append(row[1])    # password
    with open('private/command_list.csv') as csv_file:
        csv_reader = csv.reader(csv_file)
        for row in csv_reader:
            command_list.append(row[0])    # username
except:
    print("Check Private Resource!")

print(gcp_ip_list)
print(user_info)
print(command_list)

for ip in gcp_ip_list:
    print('---------------------------SSH START---------------------------')
    cli = paramiko.SSHClient()
    cli.set_missing_host_key_policy(paramiko.AutoAddPolicy)

    target = ip
    print("Target IP: ",target)
    print("USER Info.: ",user_info[0],", PW: ",user_info[1])

    cli.connect(target, port=22, username=user_info[0], password=user_info[1])
    for command in command_list:
        print("Now Command: ",command)
        stdin, stdout, stderr = cli.exec_command(command)
        #lines = stdout.readlines()
        #print(''.join(lines))
    
    cli.close()
    print('---------------------------SSH END---------------------------')