# org-scraper-lambda

## Build
```
docker build -t org-scraper-lambda . for local OR docker build --platform=linux/amd64 -t org-scraper-lambda . for deployment
```
docker tag org-scraper-lambda:latest 769325097132.dkr.ecr.eu-west-1.amazonaws.com/org-scraper-lambda:latest
```

## Push
```
aws ecr get-login-password --region eu-west-1 | docker login --username AWS --password-stdin 769325097132.dkr.ecr.eu-west-1.amazonaws.com                
docker push api-gateway-lambda:latest
```