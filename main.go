package main

import (
    "os"

    "WireguardManager/router"
)

func main() {
    router.Setup(os.Args[0])
    router.Run(os.Args[1:])
}
