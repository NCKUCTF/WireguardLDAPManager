package config

import (
    "os"
    "path/filepath"
    "github.com/joho/godotenv"
)

var WGpath string
var LDAPConf string

func init() {
    ex, err := os.Executable()
    if err == nil {
        exPath := filepath.Dir(ex)
        os.Chdir(exPath)
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
