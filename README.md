CredMan
=====
A simple and secure credential management with sync capability.

Features:
- AES encryption
- Multiple profiles
- Multiple fields per site/serivce
- Auto generate password
- Sync with custom server

# Usage
Credman works on Linux, OSX and all Windows versions.
It requires 'xsel' or 'xclip' to be installed on Linux.

First you must create a profile to store your sites in:

`credman add profile "profile name"`

It will prompt you for password. Enter a secure one and hit enter.

If it's the first profile you create, it will be set as default profile. You can add as many profiles as you want.
Commands that deal with sites require you to specify a profile name with '-p' option. If it's not specified, default profile will be used.

Next step is to add a site/service:

`credman add site "site name" --email="my email" --username="my username" --pgen`

It will create a site inside default profile with email, username and an auto generated password.
You can use `-p="profile"` to set a specific profile for site or omit to use default profile.

All credman configs and profiles will be created at user's home directory under .credman directory.

You can add/delete/rename/change password and change default profile.
Check `credman help` for a detailed help.

## Sync Servers
Credman supports sync servers to keep your passwords synced everywhere securely.
