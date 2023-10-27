package main

import (
    "os"

    "WireguardLDAPManager/router"
)

func main() {
    router.Setup(os.Args[0])
    router.Run(os.Args[1:])
}
