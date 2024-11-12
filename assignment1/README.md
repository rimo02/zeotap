# Rule Engine With AST
In this project, a rule engine that assesses circumstances based on user attributes is implemented using an Abstract Syntax Tree (AST). Using a database for storage, the engine facilitates the dynamic construction, modification, and evaluation of rules.


## Requirements
- Golang
- MongoDB
- Docker

## Steps to run the application

1.

```bash
 git clone https://github.com/rimo02/zeotap.git
```

2.

```bash
cd assignment1
```

3.

```
docker-compose up --build
```

The application should start running at port 8000
Now open your web Browser and run ``` http://localhost:8000/ ````.

OR 
1. Change `.env` file to `MONGO_URI = mongodb://localhost:27017` and run
```bash
    go run main.go
```
