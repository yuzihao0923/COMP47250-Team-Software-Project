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
- Redis
- Node.js 
- MongoDB

After installing above, build the project named "COMP47250-Team-Software-Project".
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
    - broker: "broker", "123"

3. Run consumer
```bash
cd cmd/consumer
```

```bash
go run consumer.go
```

Enter consumer's username-password pairs: "consumer", "123"

4. Run producer
```bash
cd cmd/producer
```

```bash
go run producer.go
```
Enter producer's username-password pairs: "producer", "123"

5. See the results on webpage, or check on the terminal