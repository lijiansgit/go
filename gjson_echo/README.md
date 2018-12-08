Getting Started
===============

## sample read json and echo value cmd

## vendor
* github.com/tidwall/gjson

## Installing

To start using gjson_echo, install Go and run `go get`:

```sh
$ go get -u github.com/lijiansgit/go/gjson_echo
```

## Get a value , file: test.json
```json
{
  "first"   : "cory",
  "last"    : "parker",
  "from"    : "united states",
  "age"     : 31,
  "sports"  : [ "windsurfing", "baseball", "extreeeeme kayaking" ],
  "msg"    : {
    "love": "book",
    "english_level": 8
  }
}
```


```shell
gjson_echo -f test.json -k age
gjson_echo -f test.json -k msg.love
gjson_echo -f test.json -k sports.0
```

## TODO
* ini,yaml

## json more usage please see github.com/tidwall/gjson
