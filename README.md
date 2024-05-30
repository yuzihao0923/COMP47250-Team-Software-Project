# Building a Distributed Queue System

## Team Member:

| Name         | Email                          |
|--------------|---------------------------------|
| Jingzhi Zhou | jingzhi.zhou1@ucdconnect.ie    |
| Haoyu Wang   | haoyu.wang1@ucdconnect.ie      |
| Zheng Xing   | xing.zheng@ucdconnect.ie       |
| Jiajun Zhou  | jiajun.zhou1@ucdconnect.ie     |
| Zihao Yu     | zihao.yu@ucdconnect.ie         |

## Architecture ：

<img width="100%" alt="Screen Shot 2024-05-21 at 14 21 28" src="https://github.com/yuzihao0923/COMP47250-Team-Software-Project/assets/141666207/757fa824-9865-4ff9-9a43-f9e9e99cf09d">

**Roles:**
  - **Project Manager:** Oversees project progress, ensures deadlines are met, and coordinates team activities.
  
  - **Lead Developer:** Responsible for the core implementation of the distributed queue system.
  
  - **Quality Assurance Engineer:** Ensures the system meets all quality standards through rigorous testing.
  
  - **DevOps Engineer:** Manages the cloud environment setup and deployment processes.
  
  - **Documentation Specialist:** Prepares detailed documentation for the system, including user guides and technical specifications.


## Building a Distributed Queue System Based on Redis

- Golang: [install link](https://go.dev/doc/install)  
- Redis: [install link](https://redis.io/docs/latest/operate/oss_and_stack/install/install-redis/)

---
After installing Go and Redis, build the project named "COMP47250-Team-Software-Project". 
```bash
cd COMP47250-Team-Software-Project

go mod init COMP47250-Team-Software-Project         
// This command will create "go.mod" file and 
// this file is used to manage the libraries used in this project
```

---

```
/distributed-queue-system
|-- cmd
|   |-- broker
|   |   `-- broker.go
|   |   `-- processor.go
|   |-- producer
|   |   `-- producer.go
|   `-- consumer
|       `-- consumer.go
|-- configs
|   |-- development.json
|   `-- production.json
|-- internal
|   |-- api
|   |   `-- api.go
|   |-- log
|   |   `-- log.go
|   |-- message
|   |   `-- message.go
|   |-- network
|   |   `-- network.go
|   |-- redis
|   |   `-- redis.go
|   |-- model
|   |   `-- model.go
|   `-- utils
|       `-- ip.go
|-- pkg
|   |-- config
|   |   `-- config.go
|   |-- queue
|   |   `-- queue.go
|   |-- serializer
|   |   `-- serializer.go
|   |   `-- serializerImp.go
|   `-- storage
|       `-- storage.go
|-- scripts
|   `-- deploy.sh
|-- tests
|   `-- test.go
|-- go.mod
|-- go.sum
|-- main.go
```

This cmd directory contains the code that starts the main application, usually the `main.go` file.

- **pkg:** is used to store code that can be imported by other projects.

- **internal:** is used to store code that can only be imported by this project.

- **config:** is used to store configuration files for various environments.

- **scripts:** contains scripts, such as deployment or database migration scripts.

- **tests:** contains all the test code, which may include unit tests and integration tests.


## How to Run
```bash
go run main.go   
```
Run command above on ternimal in the root folder.

