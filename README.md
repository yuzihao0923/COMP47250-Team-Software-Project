# COMP47250-Team-Software-Project

Project Title: Building a Distributed Queue System

Team Member:

Jingzhi Zhou     jingzhi.zhou1@ucdconnect.ie
Haoyu Wang       haoyu.wang1@ucdconnect.ie
Zheng Xing       xing.zheng@ucdconnect.ie
Jiajun Zhou      jiajun.zhou1@ucdconnect.ie
Zihao Yu         zihao.yu@ucdconnect.ie
| Name         | Email                          |
|--------------|---------------------------------|
| Jingzhi Zhou | jingzhi.zhou1@ucdconnect.ie    |
| Haoyu Wang   | haoyu.wang1@ucdconnect.ie      |
| Zheng Xing   | xing.zheng@ucdconnect.ie       |
| Jiajun Zhou  | jiajun.zhou1@ucdconnect.ie     |
| Zihao Yu     | zihao.yu@ucdconnect.ie         |

Architecture ：

<img width="699" alt="Screen Shot 2024-05-21 at 14 21 28" src="https://github.com/yuzihao0923/COMP47250-Team-Software-Project/assets/141666207/757fa824-9865-4ff9-9a43-f9e9e99cf09d">

Roles:

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
