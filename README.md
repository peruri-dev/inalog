# inalog

Slog based logging in Golang customized for multiple purpose

## Integration

- [x] fiber/v2

## how to install 

```
go get -v github.com/peruri-dev/inalog
```

then put in main.go

```
inalog.Init(inalog.Cfg{
    Source: true,
})

inalog.Log().Notice("hello")
inalog.Log().Info("hello, this is simple info")
inalog.Log().Error("got error!")
```

## contribution

1. Create new pull request 
2. Please follow conventional commit standard https://www.conventionalcommits.org/en/v1.0.0/
3. There will be new release after PR merged
4. So, that's why it should follow conventional commit

## release new version

1. Pull latest from main
2. Use `npx standard-version` then push the along with the tag
