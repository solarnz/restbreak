# restbreak
Small utility for running commands when the user has been active for too long

Inspired by [workrave](https://github.com/rcaelers/workrave). However this tool
is much much simplier, and instead only runs other commands when you've been
active for too long

At the moment, the tool only works under Xorg, however I plan to implement
support for Wayland Compositors that support the KDE idle protocol, such as
Sway.

Included in the repository is a sample `config.yaml` file. You will need to
copy this to `$HOME/.config/restbreak/config.yaml` before restbreak will
function.
