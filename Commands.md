# Profile management

### Adding a profile
profile/p add/a <name>

### Removing a profile
profile/p remove/rm <name>

### Renaming a profile
profile/p rename/rn <name> <new name>

### Changing profile's password
profile/p passwd/pw <name>

### Setting or getting default profile
profile/p default/d [name]

### Listing all profiles
profile/p list/l

# Site management

### Adding a new site
site/s add/a <name> [profile] [fields] [password options]

### Removing a site
site/s remove/rm <name> [profile]

### Renaming a site
site/s rename/rn <name> <new name> [profile]

### Changing or deleting a site's fields or password
site/s set/s <name> [fields] [password options] [--password] [profile]

### Listing all sites of a profile
site/s list/l [regex pattern] [profile]

### Getting all or specific fields of a site
site/s get/g <name> [fields] [profile] [--copy/-c]

# File encryption

### Encrypt a file
file/f encrypt/e <path> [name] [output file name] [--delete-original/-d] [profile]/[--no-profile/-n] [password options if -n]

### Decrypt a file
file/f decrypt/d <name/path/id> [output file] [profile]/[--no-profile/-n]

# Password restore

### Adding profile password restore
restore/rs add/a [profile]

### Restoring profile's password
restore/rs [profile]

### Removing restore for a profile
restore/rs delete/d [profile]

# Remote management

### Setting a remote
remote/r set/s <address> <username> [profile]

### Syncing with remote
remote/r sync/u [profile]

### Deleting remote
remote/r delete/d [profile]

# Server

### Starts server in background
serve

# Utility

### Generate random password
gen [--copy/-c]
