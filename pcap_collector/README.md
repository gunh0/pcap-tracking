### sshfs 

```bash
# window server
# Launch PowerShell as an administrator
Start-Service sshd

# mount
sshfs Tor@X.X.X.X:D:\\\Tor_CIFS_300 /home/jjangga94temp/Tor_CIFS -o allow_other,umask=000,reconnect,nonempty

# check
df -hT
```



### run client

```
sudo python /home/jjanggatemp/Tor_CIFS/Source/client/client.py
```