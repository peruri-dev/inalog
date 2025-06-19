# inalog

Slog based logging in Golang customized for multiple purpose

## Integration

- [x] fiber/v2

## how to install

```bash
go get -v github.com/peruri-dev/inalog
```

then put in main.go

```go
inalog.Init(inalog.Cfg{
    Source: true,
    Tinted: true,
    MessageKey: true,
})

inalog.Log().Notice("hello")
inalog.Log().Info("hello, this is simple info")
inalog.Log().Error("got error!")
```

### Configuration

Available variables for flexible configuration

```bash
INALOG_SERVICE_NAME=<your-application-name>
INALOG_SERVICE_VERSION=<your-application-version>
INALOG_SERVICE_ENV=<your-application-env>
INALOG_PRINT_PAYLOAD=<will-print-payload-headers-body-and-query-params>
INALOG_ACCESS_LOG=<will-print-every-succeeded-request>
INALOG_ERROR_LOG=<will-print-every-failed-request>
INALOG_LOG_LEVEL=<minimum-level-log-to-print> 
```

- LogLevel only accept `INFO` and `WARN` as the minimum to be printed

You can enfore payload print while `INALOG_PRINT_PAYLOAD` in false condition by injecting `_InalogForcePrint=true` in query parameters on API call

e.g.

```bash
http://localhost:9000/my-api?_InalogForcePrint=true
```

## contribution

1. Create new pull request
2. Please follow conventional commit standard <https://www.conventionalcommits.org/en/v1.0.0/>
3. There will be new release after PR merged
4. So, that's why it should follow conventional commit

## release new version

1. Pull latest from main
2. Use `npx standard-version` then push the along with the tag
