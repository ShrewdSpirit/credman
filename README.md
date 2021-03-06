Credman
=====
A simple, powerful, cross-platform and military grade (~~marketing bs~~!) secure credential management.

<p align="center"><img src="/img/demo.gif?raw=true"/></p>

## Features
- Stateless (doesn't keep passwords in memory for later use)
- Directly connect to ssh server using fields in given site (**requires cygwin or a POSIX terminal emulator on windows**)
- AES-CFB-256 and AES-CTR-256 with HMAC-SHA-256 authenticity and integrity check for profiles and files
- Standalone (single binary without any dependencies)
- Cross platform
- File encryption
- IPFS utilities and support (**WIP**)
- Multiple local profiles
- Multiple fields per site/serivce
- Auto generate a custom secure password
- Only one encrypted file per profile ~~that you can carry around!~~ so it's portable!
- Restoring profile's password in case you manage to forget it!
- Easy to use
- GraphQL interface for remotely interacting with credman daemon (**WIP**)

## Install
Binary releases are available [here](https://github.com/ShrewdSpirit/credman/releases/latest). Make sure you add the binary's directory to your PATH.

## Build from source
Requirements:
- [Go](https://golang.org/dl/) +1.11

Run the following:
```bash
git clone https://github.com/ShrewdSpirit/credman
cd credman
./build.sh install
```

## Notes
Credman works on Linux, RPi, OSX and all Windows versions (BSD **should** work, but I'll never test since I don't have the environment).
You can use it on your Android device if you have a terminal emulator (Termux is recomended).

It requires 'xsel' or 'xclip' to be installed on Linux otherwise copy function will not work.

## Basic commandline usage

First you must create a profile to store your sites in:

`$ credman profile add "profile name"`

It will prompt you for password. Enter a secure one and hit enter.

If it's the first profile you create, it will be set as default profile. You can add as many profiles as you want.
Commands that deal with sites require you to specify a profile name with '-p' option. If it's not specified, default profile will be used.

Next step is to add a site/service:

`$ credman site add "site name" --field=email="my email" --field=username="my username" --password-generate`

OR

`$ credman s a sitename -f=email='my email' -g`

It will create a site inside default profile with email, username and an auto generated password.
You can use `-p="profile"` to set a specific profile for site or omit to use default profile.
Site's fields are optional data that you can add to store extra stuff for each site.

All credman configs and profiles will be created at user's home directory under .credman directory.

**Check `credman help` or read [docs](https://github.com/ShrewdSpirit/credman/blob/master/Docs.md) for details on commands and encryption details**

## Issues
- Commandline output on windows console/ps doesn't show colors correctly ~~(blame windows)~~ so I disabled them.
- Password input doesn't work on windows git bash program idk why
