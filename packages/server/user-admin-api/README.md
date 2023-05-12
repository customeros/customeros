### Get the current social login config from Ory

```
ory get identity-config {project-id} --format yaml > identity-config.yaml
```

### Update the social login config

```
    registration:
      after:
        hooks: []
        oidc:
          hooks:
            - hook: web_hook
              config:
                url: https://user-admin-api.your-location/register
                method: POST
                body: base64://*base64-encoded-registration.jsonnet*
                response:
                  ignore: false
                  parse: false
                auth:
                  type: api_key
                  config:
                    name: X-Openline-API-KEY
                    value: ENTER_USER_ADMIN_API_KEY_HERE
                    in: header
            - hook: session
        password:
```

to build the base64 jsonnet file, use the following command:

```
cat registration.jsonnet | base64 |pbcopy
```


finally to apply the change:
```
ory update identity-config {project-id} --file identity-config.yaml
```
