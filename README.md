# go-pit

[![wercker status](https://app.wercker.com/status/80053fce485b48b7cfe2f2e9e8ba01bd/m "wercker status")](https://app.wercker.com/project/bykey/80053fce485b48b7cfe2f2e9e8ba01bd)

## SYNOPSYS

```
import(
  "github.com/naoya/go-pit"
  "log"
)

config, err := pit.Get("twitter.com")
if err != nil {
  log.Fatal(err)
}

username := config["username"]
password := config["password"]

// switch to another profile
pit.Switch("development")
config = pit.Get("twitter.com")
```

## Description

Porting [pit](https://github.com/cho45/pit) to Go

## Note

- I'm golang newbie, code review welcome
- It does not support pit.Set() yet

## TODO

- Need more documents

