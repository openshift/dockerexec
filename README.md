# dockerexec

dockerexec provides a library and a utility to execute commands in a running docker container.

## Basic usage

```
# Start a docker container
[root@localhost busybox]# docker run -it --env FOO=BAR -u daemon --rm fedora:latest bash
bash-4.3$ 
bash-4.3$ env
HOSTNAME=8b1b399c89fd
TERM=xterm
FOO=BAR
PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
PWD=/
SHLVL=1
HOME=/sbin
_=/usr/bin/env
bash-4.3$ id
uid=2(daemon) gid=2(daemon) groups=2(daemon)
bash-4.3$

# Run dockerexec specifying the same container id (shortest unique prefix works)
[root@localhost containers]# dockerexec exec  -t  --id 8b1b sh
sh-4.3$ env
HOSTNAME=8b1b399c89fd
TERM=xterm
FOO=BAR
PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
PWD=/
SHLVL=1
HOME=/sbin
_=/usr/bin/env
sh-4.3$ id
uid=2(daemon) gid=2(daemon) groups=2(daemon)
sh-4.3$ ps -ef
UID        PID  PPID  C STIME TTY          TIME CMD
daemon       1     0  0 14:31 ?        00:00:00 bash
daemon      23     0  0 14:37 ?        00:00:00 sh
daemon      28    23  0 14:37 ?        00:00:00 ps -ef
sh-4.3$ 

# Note that the user and the environment are the same for the two
# Back to the docker container
bash-4.3$ ps -ef
UID        PID  PPID  C STIME TTY          TIME CMD
daemon       1     0  0 14:31 ?        00:00:00 bash
daemon      23     0  0 14:37 ?        00:00:00 sh
daemon      29     1  0 14:37 ?        00:00:00 ps -ef
bash-4.3$
```

By default dockerexec spawns the process as the same user and env as the docker container.
User could be modified and more env variables could be specified using command line arguments.

```
[root@localhost containers]# dockerexec exec  -t --user root --env COOL=STUFF  --id 8b1b sh                         
sh-4.3# env
HOSTNAME=8b1b399c89fd
TERM=xterm
FOO=BAR
PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
PWD=/
SHLVL=1
HOME=/root
COOL=STUFF
_=/usr/bin/env
sh-4.3# id
uid=0(root) gid=0(root) groups=0(root)

[root@localhost containers]# dockerexec help exec
NAME:
   exec - execute a new command inside a container

USAGE:
   command exec [command options] [arguments...]

DESCRIPTION:
   

OPTIONS:
   --tty, -t                            allocate a TTY for the exec process
   --id                                 specify the ID of a running docker container
   --user, -u                           set the user, uid, and/or gid for the process
   --cwd                                set the current working dir
   --env '--env option --env option'    set environment variables for the process
   
```


