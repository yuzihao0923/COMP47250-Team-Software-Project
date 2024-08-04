# Building a Distributed Queue System

## Team Member:

| Name         | Email                          |Github       |
|--------------|--------------------------------|-------------|
| Jingzhi Zhou | jingzhi.zhou1@ucdconnect.ie    |kris2049     |
| Haoyu Wang   | haoyu.wang1@ucdconnect.ie      |Morgan3450   |
| Xing Zheng   | xing.zheng@ucdconnect.ie       |JettZgg      |
| Jiajun Zhou  | jiajun.zhou1@ucdconnect.ie     |JiajunZhou123|
| Zihao Yu     | zihao.yu@ucdconnect.ie         |yuzihao0923  |

## Roles:
- **Project Manager:** Jingzhi Zhou
- **Backend Developer:** Jingzhi Zhou, Zihao Yu, Xing Zheng, Haoyu Wang
- **Frontend Developer:** Jiajun Zhou, Xing Zheng
- **Testing Engineer:** Haoyu Wang, Zihao Yu
- **Documentation Specialist:** Haoyu Wang, Xing Zheng

## Architecture ：

<img width="640" alt="image" src="https://github.com/yuzihao0923/COMP47250-Team-Software-Project/assets/141666207/72510f33-cae5-4c24-a43d-105975da4988">

## Structure


  ·Project Manager: Oversees project progress, ensures deadlines are met, and coordinates team activities.
  
  ·Lead Developer: Responsible for the core implementation of the distributed queue system.
  
  ·Quality Assurance Engineer: Ensures the system meets all quality standards through rigorous testing.
  
  ·DevOps Engineer: Manages the cloud environment setup and deployment processes.
  
  ·Documentation Specialist: Prepares detailed documentation for the system, including user guides and technical specifications.


Building a Distributed Queue System Based on Redis




**********************************************************************************************
Golang : 1.22.3linux/amd64      link: https://go.dev/dl/go1.22.3.linux-amd64.tar.gz
        
Redis : redis-server 6.0.16    link: https://download.redis.io/releases/redis-6.0.16.tar.gz
**********************************************************************************************







**********************************************************************************************
After installing Go and Redis, build the project named "COMP47250-Team-Software-Project". 

~$cd COMP47250-Team-Software-Project

~/MQ$go mod init COMP47250-Team-Software-Project         // This command will create "go.mod" file and 
                            //this file is used to manage the libraries used in this project
**********************************************************************************************








**********************************************************************************************\
```
/distributed-queue-system
|-- cmd
|   |-- producer
|   |   `-- main.go
|   `-- consumer
|       `-- main.go
|-- pkg
|   |-- broker
|   |   `-- broker.go
|   |-- config
|   |   `-- config.go
|   |-- queue
|   |   `-- queue.go
|   `-- storage
|       `-- storage.go
|-- internal
|   |-- api
|   |   `-- api.go
|   `-- model
|       `-- model.go
|-- scripts
|   `-- deploy.sh
|-- configs
|   |-- development.json
|   `-- production.json
|-- tests
|   |-- integration
|   |   `-- broker_test.go
|   `-- unit
|       `-- queue_test.go
|-- go.mod
`-- go.sum
```

The cmd directory contains the code that starts the main application, usually the main.go file.

The pkg directory is used to store code that can be imported by other projects.

The internal directory is used to store code that can only be imported by this project.

The configs directory is used to store configuration files for various environments.

The scripts directory contains scripts, such as deployment or database migration scripts.

The tests directory contains all the test code, which may include unit tests and integration tests.

**********************************************************************************************

```
.
├── cmd
│   ├── broker
│   ├── consumer
│   └── producer
├── configs
│   ├── development.json
│   └── production.json
├── internal
│   ├── api
│   ├── auth
│   ├── client
│   ├── database
│   ├── log
│   ├── message
│   └── redis
├── pkg
│   ├── serializer
│   └── storage
├── scripts
│   └── deploy
├── test_data
├── tests
└── web-app
    ├── node_modules
    ├── public
    └── src
        ├── components
        ├── css
        └── services
```

- **cmd** contains the code that starts broker, consumer and producer
- **configs:** is used to store configuration files for various environments
- **internal:** is used to store code that can only be imported by this project
- **pkg:** is used to store code that can be imported by other projects
- **scripts:** contains scripts, such as deployment or database migration scripts
- **test_data:** contains fake data generation files and generation scripts
- **tests:** contains all the test code, which may include unit tests and integration tests
- **web-app:** frontend based on React

## Environment

- Golang
- Redis 6.0.16
- Node.js 22.3.0
- npm 10.8.1
- MongoDB

After installing above, build the project named "COMP47250-Team-Software-Project".
But don't forget to start MongoDB server first.

```bash
git clone https://github.com/yuzihao0923/COMP47250-Team-Software-Project.git
```

```bash
cd COMP47250-Team-Software-Project
```

```bash
go mod tidy
```

```bash
cd web-app
```

```bash
npm install
```

## How to Run locally
1. In the root folder
```bash
make start
```

2. Login broker and its username-password pairs for now are:
    - broker: "b1", "123"

3. Run consumer
```bash
cd cmd/consumer
```

```bash
go run consumer.go
```

Enter consumer's username-password pairs: "c1", "123"

4. Run producer
```bash
cd cmd/producer
```

```bash
go run producer.go
```
Enter producer's username-password pairs: "p1", "123"

5. Check results on the terminal

