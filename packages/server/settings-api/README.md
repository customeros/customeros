```bash
 curl -X POST http://localhost:10002/integration \
  -H 'Content-Type: application/json' \
  -H 'X-Openline-API-KEY: 8b010f38-e5ca-4923-a62e-9f073c5c7dbf' \
  -H 'X-Openline-USERNAME: torrey@openline.ai' \
  -d '{"hubspot": {"privateAppKey": "bob"}}'
 ```

```bash
curl -X GET http://localhost:10002/integrations \
-H 'X-Openline-API-KEY: 8b010f38-e5ca-4923-a62e-9f073c5c7dbf' \
-H 'X-Openline-USERNAME: torrey@openline.ai'
```