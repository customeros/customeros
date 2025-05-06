#docker build -t ghcr.io/customeros/customeros/neo4j-sidecar:latest -f neo4j-sidecar/Dockerfile .
#docker push ghcr.io/customeros/customeros/neo4j-sidecar:latest

docker compose -f compose.yaml pull
docker compose -f compose.yaml up -d
