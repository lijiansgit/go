# Usage
```
package main

import (
	"fmt"
	"github.com/lijiansgit/go/libs/apollo"
)

func main() {

    client := &apollo.Client{
        URL:       "http://127.0.0.1:8080",
        AppID:     "test",
        Cluster:   "default",
        Namespace: "test",
        Secret:    "",
    }

    if err := client.GetConfigCache(); err != nil {
        panic(err)
    }

    fmt.Println(client.Config)
}
```