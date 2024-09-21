# Config

How to configurate app 

## Create file

For run project create configuration file in `configs` directory.

### Local development

For local development create file with name `local.yaml`.

### Standart development

For standart development procese create file with name `config.yaml`

## File structure

File contains required values for configurations app. The configurations paramets are further described by category.

### Base configs

Base configs paramers contains information about environment app.

```yaml
env: "local"
```

### Database configs

Database paramerts contains database connections paramets and informations about database.

```yaml
database_name: "sqlite"
database:
  host: "127.0.0.1"
  port: 5432
  user: "postgres"
  pwd: "password"
  dbname: "profile"
```

### gRPC server configs

This group contains gRPC server configuration.

```yaml
grpc:
  port: 9091
```
