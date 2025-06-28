# Content-Management-System GO


## **Term of glossary**


| Term        | Meaning                   |
| ----------- | ------------------------- |
| **cms-sys** | Content-Management System |
| **lms-sys** | Learn Management System   |



----
## Folder structure of the backend
```shell

    .
├── README.md
├── cms-sys
│   ├── Dockerfile
│   ├── Makefile
│   ├── README.md
│   ├── cmd
│   │   └── main.go
│   ├── go.mod
│   ├── go.sum
│   ├── internal
│   │   ├── config
│   │   ├── handler
│   │   ├── model
│   │   ├── repository
│   │   └── service
│   ├── pkg
│   │   └── utils
│   └── test
├── docker-compose.yml
├── gateway
│   └── main.java
├── infra
│   ├── main.tf
│   ├── outputs.tf
│   └── variable.tf
├── lms-sys
│   ├── Dockerfile
│   ├── Makefile
│   ├── README.md
│   ├── cmd
│   │   └── main.go
│   ├── go.mod
│   ├── go.sum
│   ├── internal
│   │   ├── config
│   │   ├── handler
│   │   ├── model
│   │   ├── repository
│   │   └── service
│   ├── pkg
│   │   └── utils
│   └── test
└── scripts
    └── main.sh

26 directories, 19 files

```
----
# Local Development setup

Please run this command in the terminal

```shell
  # You have to change the file permission to be executable
  chmod +x ./scripts/main.sh
  
  # Running the scripts
  ./scripts/main.sh
    
```