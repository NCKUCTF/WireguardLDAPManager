package wireguard

import (
    "fmt"
    "log"
    "os"
    "os/exec"
    "bufio"
    "strings"
    "strconv"
    "path/filepath"
    "regexp"
    "net/netip"

    "github.com/google/uuid"

    "WireguardLDAPManager/utils/config"
    "WireguardLDAPManager/utils/ldap"
    "WireguardLDAPManager/utils/ipcalc"
    "WireguardLDAPManager/models/privatekey"
)

func init() {
    _, err := exec.LookPath("wg")
    if err != nil {
        log.Fatalln("Please install wireguard.")
    }
}

func GetAllName() []string {
    out, _ := exec.Command("wg", "show", "interfaces").Output()
    names := strings.Split(strings.TrimSpace(string(out)), " ")
    for i := 0; i < len(names); i++ {
        _, servervar := GetConfig(names[i])
        if _, ok := servervar["nosync"]; ok {
            names = append(names[:i], names[i+1:]...)
            i--
        }
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
        log.Fatalln("There don't have any wireguard server on this machine. Please create a wireguard server.")
    }
    fmt.Println("Please select your server name.")
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

func GetConfig(name string) (conf string, data map[string]string) {
    f, err := os.Open(filepath.Join(config.WGpath, fmt.Sprintf("%s.conf", name)))
    if err != nil {
        log.Fatalln(err)
    }
    defer f.Close()

    scanner := bufio.NewScanner(f)
    conf = ""
    data = make(map[string]string)
    ininterface := false
    rmcomment := regexp.MustCompile(`^#\s*`)
    spliteq := regexp.MustCompile(`\s*=\s*`)
    matchend := regexp.MustCompile(`^# BEGIN .*$`)
    for scanner.Scan() {
        now := scanner.Text()
        if now == "" {
            continue
        } else if now == "[Interface]" {
            ininterface = true
        } else if !ininterface {
            nowvar := rmcomment.ReplaceAllString(now, "")
            if spliteq.MatchString(nowvar) {
                data[spliteq.Split(nowvar, 2)[0]] = spliteq.Split(nowvar, 2)[1]
            }
        } else if matchend.MatchString(now) {
            break
        } else {
            if spliteq.MatchString(now) {
                data[spliteq.Split(now, 2)[0]] = spliteq.Split(now, 2)[1]
            }
        }
        conf += fmt.Sprintf("\n%s", now)
    }

    if err := scanner.Err(); err != nil {
        log.Fatalln(err)
    }

    return
}

func SetConfig(name, conf string) {
    f, err := os.OpenFile(filepath.Join(config.WGpath, fmt.Sprintf("%s.conf", name)), os.O_RDWR|os.O_TRUNC, 0600)
    if err != nil {
        log.Fatalln(err)
    }
    _, err = f.WriteString(conf)
    if err != nil {
        log.Fatalln(err)
    }
    f.Close()
    out, _ := exec.Command("wg-quick", "strip", name).Output()

    tmpname := uuid.New().String()
    f, err = os.Create(fmt.Sprintf("/tmp/%s", tmpname))
    if err != nil {
        log.Fatalln(err)
    }
    _, err = f.Write(out)
    if err != nil {
        log.Fatalln(err)
    }
    f.Close()
    defer os.Remove(fmt.Sprintf("/tmp/%s", tmpname))
    exec.Command("wg", "syncconf", name, fmt.Sprintf("/tmp/%s", tmpname)).Run()
}

func Reconfig() {
    entrys, _ := ldap.Query("", "(objectclass=wireguardKey)")
    splitcom := regexp.MustCompile(`\s*,\s*`)
    for _, name := range GetAllName() {
        conf, servervar := GetConfig(name)
        conf += "\n"
        addressesarr := splitcom.Split(servervar["Address"], -1)
        for _, entry := range entrys {
            ipindex, err := strconv.ParseInt(entry.GetAttributeValue("ipindex"), 10, 64)
            if err != nil {
                log.Fatalln(err)
            }
            addresses := []string{}
            for _, address := range addressesarr {
                nowipaddr := ipcalc.PrefixIPGet(netip.MustParsePrefix(address), ipindex)
                addresses = append(addresses, netip.PrefixFrom(nowipaddr, nowipaddr.BitLen()).String())
            }
            conf += fmt.Sprintf(
`
# BEGIN %s
[Peer]
AllowedIPs = %s
PublicKey = %s
PersistentKeepalive = %s
# END %s`, 
                fmt.Sprintf("%s-%s", config.Username, entry.GetAttributeValue("cn")),
                strings.Join(addresses, ","),
                privatekey.Pubkey(entry.GetAttributeValue("wgprivkey")),
                servervar["PersistentKeepalive"],
                fmt.Sprintf("%s-%s", config.Username, entry.GetAttributeValue("cn")),
            )
        }
        conf += "\n"
        SetConfig(name, conf)
    }
}
