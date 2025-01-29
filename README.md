# go-tado

> [!CAUTION]
> This is a side-project and still very much a work in progress.

go-tado is a Go client library for the tado° API.

## Usage

```go
import "github.com/idriesalbender/go-tado"
```

Construct a new tado° client and set up authentication using the `WithOAuthClient` method.

`WithOAuthClient` expects an initial authentication token and returns a `tado.Client` that will then automatically refresh the access token.

```go
// get intial token
token, _ := tado.DefaultOauth2Config.PasswordCredentialsToken(context.Background(), "username", "password")

// create tado client that auto-refreshes access token
client := tado.NewClient(nil).WithOAuthClient(context.Background(), nil, token)
```
