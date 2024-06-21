# Building a Distributed Queue System

## Team Member:

| Name         | Email                          |Github       |
|--------------|--------------------------------|-------------|
| Jingzhi Zhou | jingzhi.zhou1@ucdconnect.ie    |kris2049     |
| Haoyu Wang   | haoyu.wang1@ucdconnect.ie      |Morgan3450   |
| Xing Zheng   | xing.zheng@ucdconnect.ie       |JettZgg      |
| Jiajun Zhou  | jiajun.zhou1@ucdconnect.ie     |JiajunZhou123|
| Zihao Yu     | zihao.yu@ucdconnect.ie         |yuzihao0923  |

## Architecture ：

<img width="640" alt="image" src="https://github.com/yuzihao0923/COMP47250-Team-Software-Project/assets/141666207/72510f33-cae5-4c24-a43d-105975da4988">

## Roles:
- **Project Manager:** Oversees project progress, ensures deadlines are met, and coordinates team activities.
- **Lead Developer:** Responsible for the core implementation of the distributed queue system.
- **Quality Assurance Engineer:** Ensures the system meets all quality standards through rigorous testing.
- **DevOps Engineer:** Manages the cloud environment setup and deployment processes.
- **Documentation Specialist:** Prepares detailed documentation for the system, including user guides and technical specifications.

## Environment

- Golang
- Redis
- Node.js 

After installing above, build the project named "COMP47250-Team-Software-Project". 
```bash
cd COMP47250-Team-Software-Project
```

```bash
go mod init COMP47250-Team-Software-Project         
// This command will create "go.mod" file and 
// this file is used to manage the libraries used in this project
```

```bash
cd web-app
```

```bash
npm install
```

## Structure

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
│   ├── log
│   ├── message
│   ├── redis
│   └── utils
├── pkg
│   ├── config
│   ├── queue
│   ├── serializer
│   └── storage
├── scripts
│   └── deploy
├── test_data
├── tests
└── web-app
    ├── public
    └── src
        ├── components
        ├── css
        └── services
```

- **cmd** contains the code that starts broker, consumer and producer
- **pkg:** is used to store code that can be imported by other projects
- **internal:** is used to store code that can only be imported by this project
- **configs:** is used to store configuration files for various environments
- **scripts:** contains scripts, such as deployment or database migration scripts
- **tests:** contains all the test code, which may include unit tests and integration tests
- **test_data:** contains fake data generation files and generation scripts
- **web-app:** frontend based on React

## How to Run locally
1. In the root folder
```bash
make start
```

2. Login, we have 3 users, and their username-password pairs for now are:
    - broker: "broker", "123"
    - consumer: "consumer", "123"
    - producer: "producer", "123"

3. Run consumer
```bash
cd cmd/consumer
```

```bash
go run consumer.go
```

4. Run producer
```bash
cd cmd/producer
```

```bash
go run producer.go
```

5. See the results on webpage, or check on the terminal

## How to Run on Doker
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



