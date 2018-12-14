CredMan
=====
A simple and secure credential management with sync capability.

Features:
- AES encryption
- Multiple profiles
- Multiple fields per site/serivce
- Auto generate password
- Sync with custom server

# TODO
- [ ] Implement all subcommands
- [ ] Change fields to add any custom fields to sites
- [ ] Implement server sync

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
Check `credman --help` for a detailed help.

## Fields
All fields will be added to a site if specified (password field is required). Here's a complete list of fields:
- email
- username
- notes
- secqX where X is numbers from 1 to 5 for security questions

## Passwords
Credman can generate a password for each site you add or you can set your own password. The `--pgen` flag will generate a password with default settings.
You can do `credman help password` to see a list of password generator options. If you don't pass the `--pgen flag, credman will prompt you for a password.
Passwords for profiles should be entered by user and they will not be generated. If you lose your password or forget it, there will be no way to get it back!

## Sync Servers
Credman supports sync servers to keep your passwords synced everywhere securely.
