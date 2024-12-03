# EHost - Manage your `/etc/hosts` at ease

ehost is a command-line tool for easy management of your system's hosts file (/etc/hosts).

## Features

- Simple CLI interface for /etc/hosts management
- Flexible entry management (add/remove by IP, hostname, or both)
- Detailed listing capabilities with IP and hostname filtering
- Input validation for IPs and hostnames
- Built-in help system
  
## Installation

### From Source
```bash
git clone https://github.com/yourusername/ehost.git
cd ehost
make
```

## Usage

```bash
NAME:
   ehost - Manage your /etc/hosts at ease

USAGE:
   ehost [global options] command [command options]

COMMANDS:
   add      Add an entry to the hosts file
   remove   Remove entries from the hosts file
   list     List entries from the hosts file
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
```

### Command Details

#### Add Command
```bash
ehost add <ip> <hostname>
```
Adds a new entry to the hosts file. Both IP and hostname are required.

#### Remove Command
```bash
ehost remove <ip>              # Remove all entries for an IP
ehost remove <hostname>        # Remove all entries for a hostname
ehost remove <ip> <hostname>   # Remove a specific IP-hostname pair
```

#### List Command
```bash
ehost list ip <ip>            # List all hostnames for a specific IP
ehost list host <hostname>    # List all IPs for a specific hostname
```

### Examples

```bash
# Adding entries
ehost add 127.0.0.1 localhost
ehost add 192.168.1.100 myserver.local

# Removing entries
ehost remove 127.0.0.1             # Removes all entries with this IP
ehost remove myserver.local        # Removes all entries with this hostname
ehost remove 127.0.0.1 localhost   # Removes specific IP-hostname pair

# Listing entries
ehost list ip 127.0.0.1           # Shows all hostnames for 127.0.0.1
ehost list host myserver.local    # Shows all IPs for myserver.local
```
