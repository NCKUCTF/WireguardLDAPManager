package reconfig

import (
    "fmt"
    "log"
    "flag"
    "regexp"
    "strings"
    "strconv"
    "net/netip"
    "WireguardLDAPManager/models/wireguard"
    "WireguardLDAPManager/utils/ldap"
    "WireguardLDAPManager/utils/ipcalc"
)

var f *flag.FlagSet

func Usage() {
    fmt.Fprintf(f.Output(), `  %s
      Reconfig wireguard server.

Options:
  -h
        Print help message
`, f.Name())
    f.PrintDefaults()
}

func Setup(name string) {
    f = flag.NewFlagSet(name, flag.ExitOnError)
    //flag
    f.Usage = Usage
}

func Run(args []string) {
    f.Parse(args)
    subargs := f.Args()
    
    run(subargs)
}

func run(subargs []string) {
    entrys, _ := ldap.Query("", "(objectclass=wireguardKey)")
    splitcom := regexp.MustCompile(`\s*,\s*`)
    for _, name := range wireguard.GetAllName() {
        conf, servervar := wireguard.GetConfig(name)
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
# END %s`, 
                entry.GetAttributeValue("cn"),
                strings.Join(addresses, ","),
                wireguard.Pubkey(entry.GetAttributeValue("wgprivkey")),
                entry.GetAttributeValue("cn"),
            )
        }
        conf += "\n"
        wireguard.SetConfig(name, conf)
    }
}
