# Contents
<!--ts-->
 * [Profile management](#profile-management)
 * [Site management](#site-management)
 * [File encryption](#file-encryption)
 * [Password restore](#password-restore)
 * [Utility](#utility)
 * [Remote management](#remote-management)
 * [Server management](#server-management)
<!--te-->

# Profile management

### Adding a profile
`$ credman profile/p add/a <name>`

If this is the first profile, it will be set as default profile.

### Removing a profile
`$ credman profile/p remove/rm <name>`

If this is the default profile, default profile will be set to nothing after remove.

### Renaming a profile
`$ credman profile/p rename/rn <name> <new name>`

### Changing profile password
`$ credman profile/p passwd/pw <name>`

### Setting or getting default profile
`$ credman profile/p default/d [name]`

Default profile is used in management of sites, files and remotes for easier access instead of giving profile name with `--profile` option.

### Listing all profiles
`$ credman profile/p list/l`

# Site management

### Adding a new site
`$ credman site/s add/a <name> [profile] [fields] [password options] [--no-password/-n]`

### Removing a site
`$ credman site/s remove/rm <name> [profile]`

### Renaming a site
`$ credman site/s rename/rn <name> <new name> [profile]`

### Changing or deleting a site's fields or password
`$ credman site/s set/s <name> [fields] [password options] [--password] [profile]`

### Listing all sites of a profile
`$ credman site/s list/l [regex pattern] [profile]`

### Getting all or specific fields of a site
`$ credman site/s get/g <name> [fields] [profile] [--copy/-c]`

# File encryption

### Encrypt a file
`$ credman file/f encrypt/e <path> [name] [--output/-o=<output>] [--delete-original/-d] [profile]/[--no-profile/-n] [password options if -n]`

There are two ways to store your encrypted file:
1. To store it in a profile: You must specify profile (optional) and encrypted content will be stored inside profile's file. If name isn't specified, file's name with extention will be used.

2. To store it localy: The `--no-profile` flag and `--output` option and password options must be set and the encrypted file won't be stored in profile.

### Decrypt a file
`$ credman file/f decrypt/d <name/path/id> [--output/-o=<output>] [profile]/[--no-profile/-n]`

If the file isn't in profile, `--no-profile` flag must be set.

### Listing all stored files in a profile
`$ credman file/f list/l [regex pattern]`

# Password restore
This feature is disabled for all profiles by default.

Password restore might reduce the security of your profiles but in case you think you might forget your profile's password and can't rely on a pen and paper for writing down your password, you can enable this feature.

It works by asking you a few questions about your life (at least 3 questions must be answered) and encrypts your profile password by using your answers combined and hashed as key for encrypting your profile password.

The questions are defined in `<Home Directory>/.credman/.config.json` file. Be aware that changing questions might confuse you in restoring your lost password! You better not touch them or at least modify them before enabling this feature.

To reset password restore, or update the answers, you must remove restore for a profile and add it again.

### Adding profile password restore
`$ credman restore/rs add/a [profile]`

### Restoring profile's password
`$ credman restore/rs [profile]`

It starts asking you the questions you have set answers for and checks your answers against the restore key and if they are correct, your profile's password will be copied to clipboard.

### Removing restore for a profile
`$ credman restore/rs remove/r [profile]`

# Utility

### Generate random password
`$ credman gen [--copy/-c] [password options]`

# Remote management

### Setting a remote
`$ credman remote/r set/s <address> <username> [profile]`

### Syncing with remote
`$ credman remote/r sync/u [profile]`

Server will only keep your user profile file and doesn't do any magic on it. The file won't be decrypted or processed on server so all of your synced data will be safe and only the profile owner can decrypt the file. Also **KEEP IN MIND** to run a sync after you do any changes on your profile on any device you do the change. Server won't process changes or conflicts.

### Deleting remote
`$ credman remote/r remove/r [profile]`

# Server management
You can run a host for multiple clients since each client will be identified by their username and access password.

Changing a server's configuration using management commands might interrupt existing client connections.

All server management actions require you to give server password.

### Starting server
`$ credman server/sv start`

If this is the first time you start the server, you'll be prompted for server's password that will be used in server management.

After running server, the process will continue working in background and the output logs will be written at `<home directory>/.credman/logs` so you can manage your server.

### Stopping server
`$ credman server/sv stop`

### Checking server's status
`$ credman server/sv status`

### Changing server password
`$ credman server/sv passwd`

### Setting or getting listen host
`$ credman server/sv host/h [address]`

### Setting or getting listen port
`$ credman server/sv port/p [port]`

### Adding a user
`$ credman server/sv user/u add/a <username> [password options]`

### Resetting a user password
`$ credman server/sv user/u reset/r <username> [password options]`

### Deleting a user
`$ credman server/sv user/u remove/r <username>`

### Listing users
`$ credman server/sv list/l [regex pattern]`
