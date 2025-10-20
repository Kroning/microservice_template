# microservice_template

Change config.yaml and then generate service
```
go run generate.go -output ../folder-of-new-service
```

After that cd to a new directory and you can: 
1. Build image with Dockerfile. Example: 
```
docker build --network=host -t micro-template -f build/Dockerfile .
docker run docker.io/library/micro-template
```
2. Start/stop dependencies in container
```
docker compose up
docker compose down -v
```
3. Start container with service+dependencies:
```
docker compose -fdocker-compose-full.yaml up --build --force-recreate
```

Vault. If you start service localy, don't forget to add envs
```
export VAULT_ADDR=http://127.0.0.1:8200
export VAULT_ENABLE=yes
export CONFIG_PATH_APP=microservices
export CONFIG_PATH_KEY=config
```