Watchblob - WatchGuard VPN on Linux
===================================

This tiny helper tool makes it possible to use WatchGuard / Firebox / <<whatever
they are actually called>> VPNs that use multi-factor authentication on Linux.

Rather than using OpenVPN's built-in dynamic challenge/response protocol, WatchGuard
has opted for a separate implementation negotiating credentials outside of the
OpenVPN protocol, which makes it impossible to start those connections solely by
using the `openvpn` CLI and configuration files.

What this application does has been reverse-engineered from the "WatchGuard Mobile VPN
with SSL" application on OS X. A writeup of the protocol and the security implications
will be linked here in the future.

## Installation

Make sure you have Go installed and `GOPATH` configured, then simply
`go get github.com/tazjin/watchblob`.

## Usage

Right now the usage is very simple. Make sure you have the correct OpenVPN client
config ready (this is normally supplied by the WatchGuard UI) simply run:

```
watchblob vpnserver.somedomain.org username p4ssw0rd
```

The server responds with a challenge which is displayed to the user, wait until you
receive the SMS code or whatever and enter it. `watchblob` then completes the
credential negotiation and you may proceed to log in with OpenVPN using your username
and *the OTP token* (**not**  your password) as credentials.
