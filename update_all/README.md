# update-all

Running daily routines

## Preview

![update-all img preview][preview]

[preview]: ./img/preview.png

## Installation

1. build the tool with:

```sh
make
```

2. The executable binary will reside in the folder `build/main`

## How to use

**1. create a config file**

```sh
./bin/main init
```


By default, the config file will be created in `~/.config/update-all/update-all.config.yaml`

**2. override the config file**

Edit the config file, extend the list of content as you wish.

The program would read the config file and execute the commands you specified in it.

```yaml
# ~/.config/update-all/update-all.config.yaml`
- Args: [pyenv, update]
  Interval:
    Minute: 60
- Args: [pyenv, rehash]
  Interval: {}
- Args: [poetry, self, update]
  Interval:
    Hour: 24
```
