# Contents
<!--ts-->
 * [Profile management](#profile-management)
 * [Site management](#site-management)
 * [File encryption](#file-encryption)
 * [Password restore](#password-restore)
 * [Utility](#utility)
 * [Remote management](#remote-management)
 * [Server management](#server-management)
 * [Password options](#password-options)
 * [Site fields](#site-fields)
 * [Tags](#tags)
 * [Encryption details](#encryption-details)
<!--te-->

# Profile management

Each profile represents a storage for your credentials and sites. Your profile's data is encrypted with the passphrase key you provide. Credman is a credential manager that uses a password to keep your other passwords! All profiles are saved in .credman/profiles under your home or user documents directory. Never forget your profiles passwords! If you doubt your memory, you can enable [password restore](#password-restore) for your profile.

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

Sites are basically where you put the credentials for each specific service/website/etc. It might have a confusing name but you can store almost any plaintext data in sites. The sites you add to a profile will be encrypted super securely! Along with site passwords and your username, email you can add any other field that will be encrypted with the site. At this stage, you can only add text fields but in a future update, you can add text files and binary files as fields of a site.

### Adding a new site
`$ credman site/s add/a <name> [profile] [tags] [fields] [password options] [--no-password/-n]`

By default, credman adds a password field to every site you add. If you don't want any password for your site, use `--no-password` flag otherwise you must either provide [password options](#password-options) or credman will prompt you a password for your site if you already have a password. For fields, refer to [site fields](#site-fields) to see how you can use it. To set tags, refer to [tags](#tags).

### Removing a site
`$ credman site/s remove/rm <name> [profile]`

### Renaming a site
`$ credman site/s rename/rn <name> <new name> [profile]`

### Changing or deleting a site's fields or password
`$ credman site/s set/s <name> [fields] [tags] [password options] [--password] [profile]`

Refer to [password options](#password-options)

Refer to [fields](#site-fields)

Refer to [tags](#tags)

### Listing all sites of a profile
`$ credman site/s list/l [regex pattern] [tags] [profile]`

Lists all sites in a profile if pattern is not given. Otherwise you can use regex to filter sites.

### Getting all or specific fields of a site
`$ credman site/s get/g <name> [fields] [tags] [profile] [--copy/-c]`

If no fields are specified, site's password will be printed on output. If `--copy` flag is set, it will put the password in clipboard.
If more than one field is specified and `--copy` flag is used, Only first field will be copyed to clipboard (Usually password).

# File encryption

File encryption is not related to profiles and encrypted files will not be stored inside profiles. You can encrypt any kind of file with any size since the encryption/decryption is done in streaming mode.

### Encrypt a file
`$ credman file/f encrypt/e <path> [--output/-o=<output>] [--delete-original/-d]`

### Decrypt a file
`$ credman file/f decrypt/d <path> [--output/-o=<output>]`

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

Generates a random password that you can use anywhere other than credman!

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

# Password options

Some subcommands have password options that enables password generation. If you don't give password options to such subcommands, credman will prompt you to enter your own not-secure! password.

Password options are a set of commandline options that specify how the password is generated. Bellow is the full list of these options:
- `--password-generate/-g`: Enables password generator
- `--password-length/-l`: Sets the generated password length. Max length is 255 characters
- `--password-case/-c`: Sets the letter case in generated password. Valid values are `lower`/`upper`/`both`(default)
- `--password-mix/-m`: Sets the mix of characters to use in password. Valid values are `letter`/`digit`/`punc`/`all`(default)

# Site fields

Site fields are commandline options to set field values, set a list of fields to get or delete.

If you're setting or adding a new site, you must provide fields in this format: `--field=FieldKey=FieldValue` and so on.
For example to set email and username for a site you must do `-f=email=myemail@xyz.com -f=username=myniceusername`.

If you're getting a site's fields, you must provide a list of fields in this format: `--fields=Key,Key,Key` and so on.
For example to get email and username for a site you must do `-f=email,username` or `-f=email -f=username`.

Same format as above applies for deleting fields. You must provide a list of fields to `--delete` or `-d` option.

# Tags
Tags are used for grouping sites. Tags are stored per site and they use a special field that you cannot normally get or set by --field(s) option. To set tags on a site, you give your tags list to `--tags` flag. Sites that share the same tag will be grouped together in site list or they can be searched by their tags.

To delete a tag for a site, use `--delete-tags` or `-t` option for giving a list of tags to delete in `site set` subcommand.

Tags for `site get` subcommand is a flag and if set, site tags will be printed on output.

# Encryption details

#### Method
- AES-CFB-256 + HMAC SHA-256 for sites
- AES-CTR-256 + HMAC SHA-256 for files
- Scrypt for hashing

#### Salt, nounce/iv and HMAC key
Scrypt salt and HMAC key for encryption/decryption will be derivated from a master key that is generated by scrypt from the password you provide. Nounce and iv will be randomly generated.

#### Restore point
Your answers will be written to SHA-256 along with order of the answers and number of answers to generate a masterkey for encrypting your profile password.
