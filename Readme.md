# DynDns Server

[![build](https://github.com/pheelee/dyndns/actions/workflows/build.yaml/badge.svg?branch=main)](https://github.com/pheelee/dyndns/actions/workflows/build.yaml)

this is a dyndns server powered by bind9 and an api in go.
It can be used as dyndns server for e.g NAS/Router devices and acme DNS-vs01 challenges sent by the httpreq provider of traefik 

https://docs.traefik.io/https/acme/#providers

## Usage


```docker
docker run -d -p 53:53/udp -p 8081:8081 -e BASE_DOMAIN=example.com pheelee/dyndns
```

## Configuration

| Name | Example | Description |
| --- | --- | --- |
| BASE_DOMAIN | example.com | Creates a default zone for bind
| BASIC_USERNAME | dynuser | Username for txt record basic auth
| BASIC_PASSWORD | MySuperSecretPassword12 | Password for txt record basic auth

### DynDns permitted hosts

the file **dyndns.json** contains the hostnames which are permitted to update. If the file doesn't exist it will be created either using the *BASE_DOMAIN* environment variable if present or an example entry.

### Security

the endpoints `/present` and `/cleanup` can be secured by providing the environment variables **BASIC_USERNAME** and **BASIC_PASSWORD**

the endpoint `/update` is secured using the config file **dyndns.json**

```docker
docker run -d -p 53:53/udp -p 8081:8081 -e BASE_DOMAIN=example.com -e BASIC_USER=dynuser -e BASIC_PASSWORD=MySuperSecretPassword12 pheelee/dyndns
```

## Endpoints

| URL | Method(s) | Parameters | Content-Type |
| ----| :----------:| ---------- | ---------- |
| /update | GET | domain, username, password, ip4addr | ignored
| /present | POST | fqdn, value | application/json
| /cleanup | POST | fqdn, value | application/json


### Examples

* TXT Records
```bash
curl -X POST http://localhost:8081/present --user dynuser:MySuperSecretPassword12 -d '{"fqdn":"_acme-challenge.example.com","value":"1234567890"}' -H "Content-Type:application/json"
```

* A Records
```bash
curl -X GET http://localhost:8081/update?domain=dyn.example.com&username=testr&password=123&ip4addr=1.2.3.4
```
