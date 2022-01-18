# bui-api-login

## Project setup
```
go mod tidy -compat=1.17
go clean -testcache
go clean -cache -modcache -i -r
```

`.env.local`
```bash
LISTEN_ADDR=:8100

OAUTH2_REDIRECT_URL=http://localhost:8080

OAUTH2_GOOGLE_CLIENT_ID=XXXXXX.apps.googleusercontent.com
OAUTH2_GOOGLE_CLIENT_SECRET=XXXXX
OAUTH2_GOOGLE_AUTH_URL=https://accounts.google.com/o/oauth2/auth
OAUTH2_GOOGLE_TOKEN_URL=https://accounts.google.com/o/oauth2/token
OAUTH2_GOOGLE_PROFILE_URL=https://www.googleapis.com/oauth2/v1/userinfo
OAUTH2_GOOGLE_SCOPES=email

OAUTH2_FACEBOOK_CLIENT_ID=XXXXX
OAUTH2_FACEBOOK_CLIENT_SECRET=XXXXX
OAUTH2_FACEBOOK_AUTH_URL=https://www.facebook.com/v12.0/dialog/oauth
OAUTH2_FACEBOOK_TOKEN_URL=https://graph.facebook.com/v12.0/oauth/access_token
OAUTH2_FACEBOOK_PROFILE_URL=https://graph.facebook.com/v12.0/me?fields=id,name,email,picture
OAUTH2_FACEBOOK_SCOPES=email,public_profile

```

### Compiles and hot-reloads for development
```
export $(xargs <.env.local)
go run main.go
```

## TESTS
```
export T=`curl -s -X POST -d '{"email":"ggrrrr@gmail.com","password":"asdasd"}' localhost:8000/auth/login/user | jq -r '.token'`

export T=`curl -s -X POST -d '{"email":"asd@asd.com","password":"asd"}' localhost:8000/auth/login/user | jq -r '.token'`

curl -v  -H "Authorization: Bearer $T" http://localhost:8000/auth/token

```