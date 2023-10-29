package showconfig

import (
    "log"
    "fmt"
    "os"
    "flag"
    "strings"
    "strconv"
    "regexp"
    "net/netip"

    qr "github.com/skip2/go-qrcode"

    "WireguardLDAPManager/utils/config"
    "WireguardLDAPManager/utils/ldap"
    "WireguardLDAPManager/utils/ipcalc"
    "WireguardLDAPManager/models/wireguard"
    "WireguardLDAPManager/models/privatekey"
)

var f *flag.FlagSet
var qrcode bool

func Usage() {
    fmt.Fprintf(f.Output(), `  %s
  or
  %s <key name>
  or
  %s <server name> <key name>
      Show wireguard VPN config.

Options:
  -h    Print help message
`, f.Name(), f.Name(), f.Name())
    f.PrintDefaults()
}

func Setup(name string) {
    f = flag.NewFlagSet(name, flag.ExitOnError)
    f.BoolVar(&qrcode, "qrcode", false, "Display Qrcode")
    f.Usage = Usage
}

func Run(args []string) {
    f.Parse(args)
    subargs := f.Args()
    
    run(subargs)
}

func run(subargs []string) {
    keyname := ""
    servername := ""
    if len(subargs) == 1 {
        keyname = subargs[0]
    } else if len(subargs) == 2 {
        servername = subargs[0]
        keyname = subargs[1]
    } else if len(subargs) > 2 {
        fmt.Fprintln(os.Stderr, "Bad args...")
        Usage()
        return
    }
    if !wireguard.ContainName(servername) {
        if servername != "" {
            fmt.Fprintln(os.Stderr, "Server name not exist!")
        }
        servername = wireguard.ReadName()
    }
    if !privatekey.ContainName(keyname) {
        if keyname != "" {
            fmt.Fprintln(os.Stderr, "Key name not exist!")
        }
        keyname = privatekey.ReadName()
    }

    entrys, _ := ldap.Query(
        fmt.Sprintf("ou=wireguard,cn=%s,ou=people", config.Username), 
        fmt.Sprintf("(&(objectclass=wireguardKey)(cn=%s))", keyname),
    )
    _, servervar := wireguard.GetConfig(servername)
    splitcom := regexp.MustCompile(`\s*,\s*`)
    addressesarr := splitcom.Split(servervar["Address"], -1)
    ipindex, err := strconv.ParseInt(entrys[0].GetAttributeValue("ipindex"), 10, 64)
    if err != nil {
        log.Fatalln(err)
    }
    addresses := []string{}
    for _, address := range addressesarr {
        nowipaddr := ipcalc.PrefixIPGet(netip.MustParsePrefix(address), ipindex)
        addresses = append(addresses, netip.PrefixFrom(nowipaddr, netip.MustParsePrefix(address).Bits()).String())
    }
    conf := fmt.Sprintf(
`[Interface]
Address = %s
PrivateKey = %s

[Peer]
Endpoint = %s
AllowedIPs = %s
PublicKey = %s
PersistentKeepalive = %s
`, 
        strings.Join(addresses, ","),
        entrys[0].GetAttributeValue("wgprivkey"),
        fmt.Sprintf("%s:%s", servervar["Host"], servervar["ListenPort"]),
        servervar["AllowedIPs"],
        privatekey.Pubkey(servervar["PrivateKey"]),
        servervar["PersistentKeepalive"],
    )
    if qrcode {
        qrobj, err := qr.New(conf, qr.Medium)
        if err != nil {
            log.Fatalln(err)
        }
        fmt.Println(qrobj.ToSmallString(false))
    } else {
        fmt.Println(conf)
    }
}
