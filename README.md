Credman
=====
A simple, powerful, cross-platform and military grade (marketing ~~bs~~!) secure credential management.

<p align="center"><img src="/img/demo.gif?raw=true"/></p>

#### Features:
- Multiple local profiles
- Stateless (doesn't keep passwords in memory for later use)
- AES-CFB-256 and AES-CTR-256 with HMAC-SHA-256 authenticity and integrity check for profiles and files
- Identical heavy Scrypt algorithm for hashing
- Multiple fields per site/serivce
- Auto generate password
- Sync with custom server and user management
- Cross platform
- Easy usage
- Only one encrypted file per profile that you can carry around!
- Encrypt files
- Restoring profile's password in case you manage to forget it!
- Directly connect to ssh server using fields in given site (**requires cygwin or a terminal emulator on windows**)

## Installation
You can download binary releases from [here](https://github.com/ShrewdSpirit/credman/releases/latest). Make sure you put the binary somewhere that's included in your PATH environment variable.

Or you can build from source by installing [Go](https://golang.org/) and then running:

```
$ go get -u github.com/ShrewdSpirit/credman/cmd
$ go build -o $GOPATH/bin/credman github.com/ShrewdSpirit/credman/cmd
```

## Usage
Credman works on Linux, RPi, OSX and all Windows versions (BSD **should** work, but I'll never test since I don't have the environment). You can use it on your Android device if you have a terminal emulator (Termux is recomended).
It requires 'xsel' or 'xclip' to be installed on Linux otherwise copy function will not work.

First you must create a profile to store your sites in:

`$ credman profile add "profile name"`

It will prompt you for password. Enter a secure one and hit enter.

If it's the first profile you create, it will be set as default profile. You can add as many profiles as you want.
Commands that deal with sites require you to specify a profile name with '-p' option. If it's not specified, default profile will be used.

Next step is to add a site/service:

`$ credman site add "site name" --field=email="my email" --field=username="my username" --password-generate`

OR

`$ credman s a 'site name' -f=email='my email' -g`

It will create a site inside default profile with email, username and an auto generated password.
You can use `-p="profile"` to set a specific profile for site or omit to use default profile.
Site's fields are optional data that you can add to store extra stuff for each site.

All credman configs and profiles will be created at user's home directory under .credman directory.

You can add/delete/rename/change password for each site and profile.

**Check `credman help` or read [docs](https://github.com/ShrewdSpirit/credman/blob/master/Docs.md) for details on commands and encryption details**

## TODO
- [x] Implement local management
- [x] Pattern matching for site/profile names
- [x] Colorful output!
- [x] Move the functionality out of commands so other packages can use credman as a library
- [x] File encryption
- [x] Password restore
- [x] Site tags (for grouping)
- [x] SSH
- [ ] **WIP** Saving encrypted files info as sites
- [ ] Hide password field when getting site
- [ ] Use cryptographically secure random number generator
- [ ] Clear clipboard after a minute after copying a site's password to clipboard
- [ ] **WIP** Migration from older profiles
- [ ] Manually give path to profile file
- [ ] Auto update and update check
- [ ] Profile backup and auto backup
- [ ] Export to external disks/networks
- [ ] Export to gist/paste service
- [ ] Server sync
- [ ] Desktop GUI frontend
- [ ] Android app
