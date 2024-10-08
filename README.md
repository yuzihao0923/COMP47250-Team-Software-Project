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

## Code Stucture

```
.
├── cmd
│   ├── broker
│   ├── consumer
│   └── producer
│   └── proxyServer
│   └── database
├── configs
│   ├── configloader
├── internal
│   ├── api
│   ├── auth
│   ├── client
│   ├── database
│   ├── log
│   ├── message
│   └── redis
│   └── redis-cluster
├── pkg
│   ├── serializer
│   └── storage
│   ├── pool
├── scripts
│   └── deploy
├── test_data
├── tests
└── web-app
    ├── node_modules
    ├── public
    └── src
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

