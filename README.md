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
installing.

```bash
brew tap louislef299/lnet
brew install lnet
```

You can also build locally by running `make local`.

## Hardware

For a wireless adapter, I'm using the [Realtek RTL8812AU 2.4 & 5 Ghz USB Wireless Adapter][]. To install driver on Fedora, run:

```bash
sudo dnf install git dkms kernel-devel openssl
sudo git clone https://github.com/cilynx/rtl88x2bu.git \
/usr/src/rtl88x2bu-git
sudo sed -i -e "/^MAKE/s|\ssrc=\S*VERSION||
/^DEST_MODULE_LOCATION/s|/.*/net|/extra|" \
/usr/src/rtl88x2bu-git/dkms.conf
sudo dkms add rtl88x2bu/git
sudo systemctl restart dkms.service
sudo tee /etc/modules-load.d/rtl88x2bu.conf << EOF > /dev/null
88x2bu
EOF
sudo systemctl restart systemd-modules-load.service
```

[brew package manager]: https://brew.sh/
[Realtek RTL8812AU 2.4 & 5 Ghz USB Wireless Adapter]: https://zsecurity.org/product/realtek-rtl8812au-2-4-5-ghz-usb-wireless-adapter/
