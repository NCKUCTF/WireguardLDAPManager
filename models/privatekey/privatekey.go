package privatekey

import (
    "fmt"
    "log"
    "io"
    "os/exec"
    "strings"
    "strconv"

    "WireguardLDAPManager/utils/config"
    "WireguardLDAPManager/utils/ldap"
)

func GetAllName() []string {
    entrys, _ := ldap.Query(fmt.Sprintf("cn=%s,ou=people", config.Username), "(objectclass=person)")
    if len(entrys) == 0 {
        log.Fatalln("User not exist.")
    }
    entrys, _ = ldap.Query(fmt.Sprintf("ou=wireguard,cn=%s,ou=people", config.Username), "(objectclass=wireguardKey)", "cn")
    names := []string{}
    for _, entry := range entrys {
        names = append(names, entry.GetAttributeValue("cn"))
    }
    return names
}

func ContainName(name string) bool {
    for _, nowname := range GetAllName() {
        if nowname == name {
            return true
        }
    }
    return false
}

func ReadName() string {
    names := GetAllName()
    if len(names) == 0 {
        log.Fatalln("You don't have any private key. Please use `genkey <key name>` to generate a new key.")
    }
    fmt.Println("Please select your key name.")
    for i := 0; i < len(names); i++ {
        fmt.Printf("%d: %s\n", i+1, names[i])
    }
    indexstr := ""
    var index int64
    var err error
    for index, err = strconv.ParseInt(indexstr, 10, 64); err != nil || int(index) < 1 || int(index) > len(names); index, err = strconv.ParseInt(indexstr, 10, 64) {
        fmt.Printf("> ")
        fmt.Scanln(&indexstr)
    }
    return names[index - 1]
}

func ReadNewName() string {
    fmt.Println("Please input your new key name.")
    name := ""
    fmt.Printf("> ")
    fmt.Scanln(&name)
    return name
}

func Pubkey(privkey string) string {
    cmd := exec.Command("wg", "pubkey")
    stdin, err := cmd.StdinPipe()
    if err != nil {
        log.Fatalln(err)
    }
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        log.Fatalln(err)
    }
    io.WriteString(stdin, privkey)
    stdin.Close()
    if err = cmd.Start(); err != nil {
        log.Fatalln(err)
    }
    out, _ := io.ReadAll(stdout)
    if err = cmd.Wait(); err != nil {
        log.Fatalln(err)
    }
    return strings.TrimSpace(string(out))
}

func Generate(keyname string) {
    err := exec.Command("ldapwgkeyadd", "-f", config.LDAPConf, keyname, config.Username).Run()
    if err != nil {
        log.Fatalln(err)
    }
}

func Delete(keynames []string) {
    err := exec.Command("ldapwgkeydel", "-f", config.LDAPConf, "-k", strings.Join(keynames, ","), config.Username).Run()
    if err != nil {
        log.Fatalln(err)
    }
}

func Clear() {
    err := exec.Command("ldapwgkeydel", "-f", config.LDAPConf, "-c", config.Username).Run()
    if err != nil {
        log.Fatalln(err)
    }
}

