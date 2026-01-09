# VPS SETUP

To test the script locally you can use Ubuntu's mutlipass.

**Install it**
```
sudo snap install multipass
```

**Create the VM (Ubuntu 24.04 is 'noble')**
```
multipass launch noble --name test-bench
```

**Transfer the script to the VM**
```
multipass transfer bin/vps-setup test-bench:/home/ubuntu/
```

**Run the script via SSH**
```
multipass exec test-bench -- sudo /home/ubuntu/vps-setup
```

**Check results, then Burn it down**
```
multipass delete test-bench
multipass purge
```