# HLK-SW16
[![Build Status](https://www.travis-ci.org/vitaly-kashtalyan/hlk-sw16.svg?branch=master)](https://www.travis-ci.org/vitaly-kashtalyan/hlk-sw16)

## Description:
Package hlk-sw16 created to help manage 16x relay switches electronic module.

### Install package:
``` bash
go get -u github.com/vitaly-kashtalyan/hlk-sw16
```
You can also manually git clone the repository to:
``` bash
$GOPATH/src/github.com/vitaly-kashtalyan/hlk-sw16
```

### Usage:
``` go
package main

import "github.com/vitaly-kashtalyan/hlk-sw16"

func main() {
	hlk := hlk_sw16.New("192.168.16.254", "8080")
	if hlk.Err == nil {
	    //
		_ = hlk.StatusRelays()
		msg, _ := hlk.ReadMessage()
		fmt.Println("Message: ", msg)
		_ = hlk.Close()
	}
}
```