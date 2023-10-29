package config

import (
    "log"
    "fmt"
    "strings"
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

func Ask(message string) bool {
    data := ""
    for strings.ToLower(data) != "y" && strings.ToLower(data) != "n" {
        fmt.Printf("%s [y/n]: ", message)
        fmt.Scanln(&data)
    }
    return strings.ToLower(data) == "y"
}
