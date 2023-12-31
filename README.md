# WireguardLDAPManager

Append your wireguard peer by run this command in your ldap account.

## Build
1. Install golang 1.19
2. Clone this repo and cd into WireguardLDAPManager
``` bash
git clone https://github.com/NCKUCTF/WireguardLDAPManager
cd WireguardLDAPManager
```
3. Run make
``` bash
make clean && make
```

## Install
1. Install [LDAPUserTools](https://github.com/Jimmy01240397/LDAPUserTools) or [LDAPCLI](https://github.com/NCKUCTF/LDAPCLI) and set your bind config in `/etc/ldap/bind.yaml`
```
vi /etc/ldap/bind.yaml
```
```yaml
url: ldap://127.0.0.1
binddn: cn=admin,dc=example,dc=com
bindpasswd: test1234
```
2. CD into repo dir
``` bash
cd WireguardLDAPManager
```
3. Run make
```
make install
```

## Uninstall
1. CD into repo dir
``` bash
cd WireguardLDAPManager
```
2. Run make
```
make uninstall
```

## Run
Use `wgldapmanager help` to see command usage.

```
Usage: wgldapmanager COMMAND

Options:
  -h    Print help message

A list of commands is shown below. To get detailed usage and help for a
command, run:
  wgldapmanager help COMMAND

Here is the list of commands available with a short syntax reminder. Use the
'help' command above to get full usage details.

  help
  reconfig
  showconfig <key name>
  genkey <key name>
  delkey <key name>
  clearkey
```
