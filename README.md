
# i3-icon-to-go

- Renames i3wm workspaces according to the workspace window names
- Allows usage of [Font Awesome](https://origin.fontawesome.com/icons?d=gallery) icons instead of app names

## Setup
1. To support icons Font Awesome should be available on your system. In case it is not installed - use package manager to install it.
You can use `fc-list | grep Awesome` or just `i3-icon-to-go awesome` to check the Font Awesome availability

2. **configs** directory contains sample configuration files in yaml format
These files should be placed either under `~/.i3` or `~/.config/i3` directory.
`fa-icons.yaml` sets one-to-one mapping from icon name to UTF-8 code as set by Font Awesome.
`app-icons.yaml` sets one-to-many mapping from icon name to app name
A default `fa-icons.yaml` can be produced by executing `i3-icon-to-go parse > ~/.config/i3/fa-icons.yaml`

3. Just place the executable file anywhere and add this line to your i3 config:
`exec_always --no-startup-id i3-icon-to-go` 

## Command line parameters
```
  -c         path to the app-icons.yaml config file
  -u         display only unique icons. Default is True
  -l         trim app names to this length. Default is 12
  -d         app delimiter. Default is a pipe character "|"
```

Inspired by https://github.com/cboddy/i3-workspace-names-daemon
