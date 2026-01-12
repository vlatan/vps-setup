# VPS SETUP

The automated remote machine configuration is done via a Golang script. The script is tailored to run on Ubuntu 24.04 LTS. First test it locally before using it to setup a remote production machine. This will give you information on what answers and/or variables should you prepare beforhand.


## Testing the Script Locally

To test the script locally you can use Ubuntu's [mutlipass](https://github.com/canonical/multipass).

Install mutlipass.
```
sudo snap install multipass
```

Create the VM (Ubuntu 24.04 is 'noble').
```
multipass launch noble --name test-bench
```

Compile and transfer the binary to the VM.
```
make
multipass transfer bin/vps-setup test-bench:/home/ubuntu/
```

Run the script via SSH.
```
multipass exec test-bench -- sudo /home/ubuntu/vps-setup
```

Or entrer the VM's shell and do work inside.
```
multipass shell test-bench
```

Burn it down.
```
multipass delete test-bench
multipass purge
```


## Run the Script on a Remote Machine

First of course you need to create your VPS at your cloud provider.

Compile the binary, move it to your VPS, login, run it and remove it.
```
make
scp bin/vps-setup root@xx.xxx.xxxx.xxxx:
ssh root@xx.xxx.xxxx.xxxx
./vps-setup
rm vps-setup
```

The script will guide you trough the process, ask you to provide input. It will install all the necessary software and configure the machine.


## Post Instalation Steps - Configure SSH Key-Based Authentication

In order to make the login process to your remote machine more streamlined we need to make some changes


#### Configure SSH Locally

**ON YOUR LOCAL MACHINE** edit your local `~/.ssh/config` file and add this content to it.
```
Host <remote_host>
Hostname xxx.xx.xx.xx
Port <port>
User <remote_username>
PubKeyAuthentication yes
IdentityFile /home/<local_username>/.ssh/<key_name>
```

Restrict access to config file.
```
chmod 600 ~/.ssh/config
```

Generate private and public key. When asked supply the path where you want them to be saved which will be `/home/<local_username>/.ssh/<key_name>`.
```
ssh-keygen -t ed25519 -C "<your_comment>"
```

Restrict access to the keys.
```
chmod 400 ~/.ssh/<key_name>.pub ~/.ssh/<key_name>
```

Push the public key to the remote server. This will still ask for user password.
```
ssh-copy-id -i ~/.ssh/<key_name>.pub -p <port> <remote_username>@xxx.xx.xxx.xx
```

If you are able to connect with the command `ssh <remote_host>` then you can procede to finish the SSH hardening of the remote machine.


#### Disabe Password Authentication and Root Login on the Remote Server

Login to the **REMOTE** machine, open the `/etc/ssh/sshd_config.d/harden.conf` file and comment out the rules inside to disable the root login altogether, restart the `ssh` service for changes to take effect. You can also safely upgrade the software and reboot.
```
ssh <remote_host>
sudo nano /etc/ssh/sshd_config.d/harden.conf
sudo systemctl daemon-reload
sudo systemctl restart ssh
sysupdate
sudo reboot
```
