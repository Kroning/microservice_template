# microservice_template

Change config.yaml and then generate service
```
go run generate.go -output ../folder-of-new-service
```

After that cd to a new directory and you can: 
1. Build image with Dockerfile (if you need). Example: 
```
docker build --network=host -t micro-template -f build/Dockerfile .
docker run docker.io/library/micro-template
```
2. Start/stop only dependencies in container:
```
docker compose -fbuild/docker-compose.yaml up
docker compose down -v
```
3. Start container with service+dependencies:
```
docker compose -fbuild/docker-compose-full.yaml up --build --force-recreate
```

Vault. If you start service locally, don't forget to add envs
```
export VAULT_ADDR=http://127.0.0.1:8200
export VAULT_ENABLE=yes
export CONFIG_PATH_APP=microservices
export CONFIG_PATH_KEY=config
```