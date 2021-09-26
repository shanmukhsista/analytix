```shell
docker run  --memory=4g --cpus=2 -p 9000:9000 -p 8123:8123 --ulimit nofile=262144:262144 yandex/clickhouse-server
```