# terminal-2fa

A sample 2FA tools in terminal.

## Install

* Use go command.

```shell
go install github.com/dawenxi-tech/terminal-2fa@latest
```

Note: 2fa will use `$HOME/.config/2fa/config.json` to store secrets.
The secret in `config.json` was AES encrypt by fixed key and If you want
to use a dynamic key to encrypt secret, try to use next command to install

```shell
go install -ldflags "-X main.encryptKey=`head -c 16 /dev/urandom|base64`"  github.com/dawenxi-tech/terminal-2fa@latest
```

*IMPORTANT*: this will use a dynamic and unique key to encrypt secret thus if you uninstall `2fa`, it will never can
find back your data.

## Usage

* Display 2fa

```shell
2fa
```

* Config by command line

```shell
# add new 2fa
2fa config add -name foo -secret HDDLTQ2TIMDX24PU

# edit 2fa 
2fa c edit -id 1 -name bar

# delete 2fa
2fa c delete -id 1

# import 2fa form other app export url
2fa c import -url "otpauth-migration://offline?data=CjEKCkhlbGxvId6tvu8SGEV4YW1wbGU6YWxpY2VAZ29vZ2xlLmNvbRoHRXhhbXBsZTAC"

# move up
2fa c move -id 1 -offset -1

# move down
2fa c move -id 1 -offset 1
```

*Config by terminal ui

```shell
2fa gui
```


