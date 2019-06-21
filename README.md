Credman
=====
A simple, powerful, cross-platform and military grade (marketing ~~bs~~!) secure credential management.

<p align="center"><img src="/img/demo.gif?raw=true"/></p>

#### Features:
- Easy usage!
- Multiple local profiles
- Stateless (doesn't keep passwords in memory for later use)
- AES-CFB-256 and AES-CTR-256 with HMAC-SHA-256 authenticity and integrity check for profiles and files
- Identical scrypt password hashing
- Multiple fields per site/serivce
- CLI and GUI (GUI is in WIP!)
- Standalone!
- Auto generate password
- Sync with custom server and user management
- Cross platform
- Only one encrypted file per profile that you can carry around!
- File encryption
- Restoring profile's password in case you manage to forget it!
- Directly connect to ssh server using fields in given site (**requires cygwin or a POSIX terminal emulator on windows**)

## Install
Binary releases are available [here](https://github.com/ShrewdSpirit/credman/releases/latest). Make sure you add the binary's directory to your PATH.

Darwin builds don't come with GUI since I can't cross-compile for mac with cgo enabled. If you're a darwin user and you want GUI, you must build credman from source.

## Build from source
Requirements:
- Go toolchain
- GCC
- nodejs and npm/yarn

Run the following:
```bash
go get -u github.com/ShrewdSpirit/credman/cmd/credman/...
go install github.com/ShrewdSpirit/gassets/cmd/gassets
cd $GOPATH/src/github.com/ShrewdSpirit/credman/gui/assets
npm install
npm run build
gassets -d .
cd ../../cmd/credman
go install -tags='gui' # omit the -tags='gui' if you don't want GUI
```

## CLI Usage
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
- [x] File encryption
- [x] Password restore
- [x] Site tags (for grouping)
- [x] SSH
- [x] Saving encrypted files info as sites
- [x] Hide password field when getting site
- [x] Use cryptographically secure random number generator
- [x] Manually give path to profile file
- [x] Clear clipboard after 10 seconds for copying a site's password
- [x] Auto update and update check
- [x] Desktop GUI
- [ ] Server sync

## Issues

- Commandline output on windows console/ps doesn't show colors correctly (blame windows)
- Password input doesn't work on windows git bash program
