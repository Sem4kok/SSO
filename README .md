
# 🔒 Single-Sign On [SSO] 🔒

🧩 🔬 An SSO service that supports a microservice architecture {auth, user-info, premission} - microservices. 

🛡️ ⚔️ Responsible for registering users, storing their data in a secure form in the database, authentication, checking for administrator rights.

🌐 🛠️ The server uses Remote Procedure Call is a technique for building distributed systems.


## Technologies Stack

- jwt tokens-[JWT](https://github.com/golang-jwt/jwt) ㅤㅤㅤㅤㅤㅤㅤ||  crypto - [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)
- env\flag parse-[cleanenv](https://github.com/ilyakaznacheev/cleanenv)  ㅤㅤㅤ|| logger - [slog](https://pkg.go.dev/log/slog) 
- request generator -[gofakeit](https://github.com/brianvoe/gofakeit)  ㅤ|| testing - [testify](https://github.com/stretchr/testify)
- db migration-[Goose](https://github.com/pressly/goose) ㅤㅤㅤㅤㅤ|| db - sqlite
- rpc framework-[gRPC](https://grpc.io/)ㅤㅤㅤㅤㅤ|| client - [postman](https://www.postman.com/)
## API Reference

### Register

```rpc
  Auth/Register
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `email` | `string` | **Required**. User email|
| `password` | `string` | **Required**. User password|

```JSON
{
    "email": "Goblin@sosamuzik.com",
    "password": "FreeSosaPlatina67@!"
}
```

### Login

```rpc
  Auth/Login
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `email`      | `string` | **Required**. User email|
| `password`      | `string` | **Required**. User password |
| `app_id`      | `int32` | User has to choose app |

```JSON
{
    "email": "Goblin@sosamuzik.com",
    "password": "FreeSosaPlatina67@!"
    "app_id": 1
}
```

### IsAdmin

```rpc
  Auth/IsAdmin
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `user_id`      | `int64` | **Required**.|

```JSON
{
    "app_id": 8952812
}
```
## Run Locally

Clone the project

```bash
  git clone https://github.com/Sem4kok/SSO
```

Go to the project directory

```bash
  cd SSO
```

Make DataBase migrations

```bash
  make migrate
```

Start the auth server

```bash
  make auth
```

If you want to rewrite protobuf contract then you have to 
```bash
  make protoc
```

