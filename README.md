# VPS SETUP

The automated remote machine configuration is done via a Golang script. The script is tailored to run on Ubuntu 24.04 LTS. First you might want to test it locally before using it to setup a remote production machine.


## Prerequisites

Fill out an `.env` file like it's shown in [example.env](example.env).

**Note**: Don't forget to generate an SSH keys pair on your **local machine**.
``` bash
ssh-keygen -t ed25519 -f /home/<local_username>/.ssh/<key_name> -C "<your_comment>"
```

Restrict access to the keys.
``` bash
chmod 400 ~/.ssh/<key_name>.pub ~/.ssh/<key_name>
```

The SSH public key should be inside the `~/.ssh/<key_name>.pub` file. The string should be in the format `[type] [key] [comment]`. Supply this key to the `.env` file.


## Testing the Script Locally

To test the script locally you can use Ubuntu's [mutlipass](https://github.com/canonical/multipass).

``` bash
# Install mutlipass
sudo snap install multipass
```

Compile the binary, launch a multipass instance, move the `.env` file and the binary to the instance, note the instance's IP, execute the script and try to login.

``` bash
make # Compile the binary

# Create the VM (Ubuntu 24.04 is 'noble')
multipass launch noble --name test-bench

# Transfer the .env file and the binary to the VM
multipass transfer .env bin/vps-setup test-bench:/home/ubuntu/

# Before running the script note the test-bench instance IP
multipass list

# If you want to look around inside the instance
# before running the script
multipass shell test-bench

# Execute the script via SSH.
# This will change the SSH port of the instance,
# and enable login JUST for your user with key pairs
multipass exec test-bench -- sudo /home/ubuntu/vps-setup

# Login to the instance.
# Replace the user, ip, port and the pubkey.
ssh user@ip -p port -i ~/.ssh/pubkey

# Burn it down
multipass stop --force test-bench &&
multipass delete --purge test-bench
```


## Run the Script on a Remote Machine



Compile the binary, move the `.env` file and the binary to your VPS, login, execute the binary, remove the files and reboot.
``` bash
make
scp .env bin/vps-setup root@xx.xxx.xxxx.xxxx:
ssh root@xx.xxx.xxxx.xxxx
./vps-setup
rm .env vps-setup
reboot
```

After the reboot you should be able to login to you VPS by using:
``` bash
ssh -p <port> <remote_username>@xx.xxx.xxxx.xxxx
```


### Post Instalation Steps 

In order to make the login process to your remote machine more streamlined you need to make some minor changes on your **local machine**.


Edit your local `~/.ssh/config` file and append this content to it.
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

Now you should be able to login to your VPS by simply using `ssh <remote_host>`.
