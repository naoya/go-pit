# go-pit

## SYNOPSYS

```
import(
  "github.com/naoya/go-pit"
)

config := pit.Get("twitter.com")

username := config["username"]
password := config["password"]

// switch to another profile
pit.Switch("development")
config = pit.Get("twitter.com")
```

## Description

Porting [pit](https://github.com/cho45/pit) to Go

## Note

- It does not support pit.Set() yet
