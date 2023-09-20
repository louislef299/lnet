# Louis' Network Tool

This is a network cli tool that is in the process of getting fully fleshed out.
Feel free to play around with it, it's super basic and not that hard to get the
hang of. It currently just runs a port scan and a couple dns queries. Hopefully
will get around to adding more capabilities as this would be nice to use for the
average system administrator.

Check out the [command docs](docs/cmds) for more command-specific information.

## Releasing

This repo uses release please along with go releaser in order to automatically
produce artifacts in GitHub

## Installation

You can build locally by running `make local` if you don't use the [brew package
manager][]. With homebrew, you can add lnet to your system by first tapping, then
installing. Currently supports linux, windows and darwin(macOS).

```bash
brew tap louislef299/lnet
brew install lnet
```

You can also build locally by running `make local`.

## Hardware

For a wireless adapter, I'm using the [Realtek RTL8812AU 2.4 & 5 Ghz USB
Wireless Adapter][]. To install driver on Fedora, I am using the [public git
repo][]. To check to see if the driver is successfully installed, run `inxi
-Nxx`.

```bash
ifconfig wlan0 down
ifconfig wlan hw ether 00:22:33:44:55:66
ifconfig wlan0 up
```

## Issues with `sudo`

If you are getting issues running `sudo lnet` with `command not found`, a simple
way to get around this is to make a symbolic link in `usr/bin` to your go bin(if
that is where you have it installed).

Saving this link here because it looks like there is someone who is on a similar
journey to me: [mdlayher][].

[brew package manager]: https://brew.sh/
[public git repo]: https://github.com/gnab/rtl8812au
[Realtek RTL8812AU 2.4 & 5 Ghz USB Wireless Adapter]: https://zsecurity.org/product/realtek-rtl8812au-2-4-5-ghz-usb-wireless-adapter/
[mdlayher]:(https://mdlayher.com/blog/linux-netlink-and-go-part-1-netlink/)
