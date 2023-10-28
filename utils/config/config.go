package config

import (
    "log"
    "syscall"
    "os"
    "os/user"
    "path/filepath"
    "github.com/joho/godotenv"
)

var WGpath string
var LDAPConf string
var Username string

func init() {
    nowuser, err := user.Current()
    Username = nowuser.Username
    ex, err := os.Executable()
    if err == nil {
        exPath := filepath.Dir(ex)
        os.Chdir(exPath)
    }
    err = syscall.Setuid(syscall.Geteuid())
    if err != nil {
        log.Fatalln(err)
    }
    err = godotenv.Load()
    exists := false
    WGpath, exists = os.LookupEnv("WGPATH")
    if !exists {
        WGpath = "/etc/wireguard"
    }
    LDAPConf, exists = os.LookupEnv("LDAPCONF")
    if !exists {
        LDAPConf = "/etc/ldap/bind.yaml"
    }
}
