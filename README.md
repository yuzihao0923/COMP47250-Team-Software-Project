# Building a Distributed Queue System

## Team Member:

| Name         | Email                          |
|--------------|---------------------------------|
| Jingzhi Zhou | jingzhi.zhou1@ucdconnect.ie    |
| Haoyu Wang   | haoyu.wang1@ucdconnect.ie      |
| Xing Zheng   | xing.zheng@ucdconnect.ie       |
| Jiajun Zhou  | jiajun.zhou1@ucdconnect.ie     |
| Zihao Yu     | zihao.yu@ucdconnect.ie         |

## Architecture ：

<img width="640" alt="image" src="https://github.com/yuzihao0923/COMP47250-Team-Software-Project/assets/141666207/72510f33-cae5-4c24-a43d-105975da4988">

## Roles:
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
|   |-- redis
|   |   `-- redis.go
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
|-- test_data
|   `-- main.go
|-- go.mod
|-- go.sum
|-- main.go
```

This cmd directory contains the code that starts the main application, usually the `main.go` file.

- **pkg:** is used to store code that can be imported by other projects.

- **internal:** is used to store code that can only be imported by this project.

- **configs:** is used to store configuration files for various environments.

- **scripts:** contains scripts, such as deployment or database migration scripts.

- **tests:** contains all the test code, which may include unit tests and integration tests.
- **test_data:** contains fake data generation files and generation scripts


## How to Run
```bash
docker-compose up --build 
```
Run command above on ternimal in the root folder.


docker network create redis --subnet 172.38.0.0/16
cluster nodes
- a7a44b8161dadf23871254bfabeb9bfc5b3870e1 172.38.0.12:6379@16379 master - 0 1718564344000 2 connected 5461-10922
- 4ab16e0ee8bfd4f0e7cafdfbbbee86a4846dca16 172.38.0.13:6379@16379 master - 0 1718564344537 3 connected 10923-16383
- 49502ad45e6b877a04ef77d72545504045c3684a 172.38.0.15:6379@16379 slave - 20b1e40d5f234b5d41012acc21ff896a9761f04b 0 1718564344436 5 connected
- 238f7ae0da1362e6b7d8d38e02fafbee0de62389 172.38.0.14:6379@16379 slave 4ab16e0ee8bfd4f0e7cafdfbbbee86a4846dca16 0 1718564343433 4 connected
- 20b1e40d5f234b5d41012acc21ff896a9761f04b 172.38.0.11:6379@16379 myself,master - 0 1718564342000 1 connected 0-5460
- e381dd355e98a297661a4d31cc9db8baf20cbca1 172.38.0.16:6379@16379 slave a7a44b8161dadf23871254bfabeb9bfc5b3870e1 0 1718564342429 6 connected


