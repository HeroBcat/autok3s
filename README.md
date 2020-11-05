# autok3s

- Install docker.io k3s nfs-common automatically, and the commands you pass in
- Import authorized_keys.txt to `.ssh/authorized_keys` with username and password for ssh login
 
> autok3s --user ubuntu --password password -m master_ip -w worker_ip1 -w worker_ip2 --auth_keys_path /path/authorized_keys.txt \\<br/>
> 	--master_extra_args "--docker" \\<br/>
> 	--worker_extra_args "--docker" \\<br/>
> 	--master_commands "sudo apt-get install nfs-kernel-server" \\<br/>
> 	--master_commands "sudo mkdir -p /nfs" \\<br/>
> 	--master_commands "sudo touch /etc/exports" \\<br/>
> 	--master_commands \\""sudo echo '/nfs *(rw,sync,no_subtree_check,no_root_squash)' >> /etc/exports"\\" \\<br/>
> 	--master_commands "sudo exportfs -ar" \\<br/>
> 	--master_commands "sudo /etc/init.d/nfs-kernel-server restart" \\<br/>
> 	\# --pre_master_commands \\<br/>
> 	\# --pre_master_commands \\<br/>
> 	\# --pre_master_commands \\<br/>
> 	\# --pre_worker_commands \\<br/>
> 	\# --pre_worker_commands \\<br/>
> 	\# --pre_worker_commands \\<br/>
> 	\# --worker_commands \\<br/>
> 	\# --worker_commands \\<br/>
> 	\# --worker_commands \\
