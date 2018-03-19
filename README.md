# Zero Bundle
runs a mini env flist on linux.

## Operation
zbundle, uses 0-fs to mount an flist then automatically chroot, and run the start point of the env
- once the flist is mounted, zbundle reads the `/etc/env` file
- next, zbundle runs `/etc/start` inside the chroot with the env variables loaded from `/etc/env`
- zbundle shows all the output of the start process on stdout, and stderr and exits once the start process exits

## `env` file format
since NOT all flists will have a shell (`bash`, `sh`, or others) it will not be possible to initialize the
env vars of the `start` process. A special (optional) env file is used to setup the env vars of the start process.

> Note: all env vars from the `env` file can be overridedn from the `zbundle` command line with the `-e` flag

example env file
```bash
# set up the path
PATH=/bin:/usr/bin

# set up the version
VERSION=1.0
```

## `start` file
After loading the `env` from the `/etc/env` file, next `zbundle` will execute the `/etc/start` script.
All the env variables from the `env` file, plus the ones defined from the command line are gonna be available
to the `start` excutable.

- `/etc/start` MUST have `excute` permission (`chmod a+x /etc/start`)
- if `start` is a shell script, the excuting `shebang` must actually be part of the flist otherwise it will not work, because the flist
runs in a chroot, there is noway to access the binaries from the host.

### Notes
- the env file is NOT a shell script
- empty lines and lines starting with `#` are ignored

## Example
```
zbundle https://hub.gig.tech/azmy/bundle-test.flist mount-location
```

## Help

```
~ zbundle --help
NAME:
   zbundle - run mini environments from flist

USAGE:
   0-bundle [options] flist root

VERSION:
   1.0

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --env value, -e value      custom environemt variables
   --redis value, -r value    Redis server address for error reporting
   --storage value, -s value  storage url to use (default: "ardb://hub.gig.tech:16379")
   --debug, -d                run in debug mode
   --no-exit                  do not terminate (unmount the sandbox) after /etc/start exits
   --help, -h                 show help
   --version, -v              print the version
   ```