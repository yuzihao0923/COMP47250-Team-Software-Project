from locust import HttpUser, TaskSet, task, between
import requests
import json
import time
import random
import string
import os

proxy_url = "http://localhost:8888"
user_consumer = {"username": "c1", "password": "123"}

class ConsumerTasks(TaskSet):
    def on_start(self):
        self.broker_addr = self.get_broker_address()
        self.token = self.authenticate_user(user_consumer)
        # self.clear_json_file()
        self.stream_name = ''.join(random.choices(string.ascii_letters, k=10))
        self.group_name = ''.join(random.choices(string.ascii_letters, k=10))
        self.write_to_json_file()
        self.register()

    def clear_json_file(self):
        file_path = "../message/messages.json"
        # Ensure the directory exists
        os.makedirs(os.path.dirname(file_path), exist_ok=True)
        # Clear the file
        with open(file_path, "w") as file:
            json.dump([], file)
    
    def read_from_json_file(self):
        try:
            with open("../message/messages.json", "r") as file:
                data = json.load(file)
                return data
        except (IOError, json.JSONDecodeError) as e:
            print(f"Error reading JSON file: {e}")
            return []
    def write_to_json_file(self):
        message_data = {
            "stream_name": self.stream_name,
            "group_name": self.group_name
        }
        
        file_path = "../message/messages.json"
        
        # Ensure the directory exists
        os.makedirs(os.path.dirname(file_path), exist_ok=True)
        
        # Read existing data or initialize an empty list
        try:
            with open(file_path, "r") as file:
                data = json.load(file)
        except (IOError, json.JSONDecodeError):
            data = []
        
        data.append(message_data)
        
        # Write updated data back to the file
        with open(file_path, "w") as file:
            json.dump(data, file, indent=4)

    def get_broker_address(self):
        response = requests.get(f"{proxy_url}/get-broker")
        response.raise_for_status()
        address = response.json()["address"]
        return address

    def authenticate_user(self, user):
        broker_login_url = f"http://{self.broker_addr}/login"
        response = self.client.post(broker_login_url, json={"username": user["username"], "password": user["password"]})
        response.raise_for_status()
        data = response.json()
        if data["role"] == "producer" or data["role"] == "consumer":
            return data["token"]
        else:
            raise Exception(f"User {user['username']} is not authorized")
    

    def register(self):
        msg = {
            "type": "registration",
            "consumer_info": {
                "stream_name": self.stream_name,
                "group_name": self.group_name
            }
        }
        headers = {"Authorization": f"Bearer {self.token}", "Content-Type": "application/json"}
        response = self.client.post(f"http://{self.broker_addr}/register", json=msg, headers=headers)
        if response.status_code == 200:
            print(f"Consumer registered: {response.text}")
        else:
            print(f"Error registering consumer: {response.text}")


    def ack(self,msg):
        headers = {"Authorization": f"Bearer {self.token}", "Content-Type": "application/json"}
        response = self.client.post(f"http://{self.broker_addr}/ack", json=msg, headers=headers)
        if response.status_code == 200:
            print(f": {response.text}")
        else:
            print(f"Error registering consumer: {response.text}")

    @task
    def consume_messages(self):
        # 从JSON文件中读取消息流信息列表
        nameInfoList = self.read_from_json_file()

        while True:  # 无限循环，依次从每个stream消费消息
            for nameInfo in nameInfoList:
                # print("---------------" + "Consume from "+str(nameInfo["stream_name"]) + "------------")

                try:
                    # 发送GET请求以消费消息
                    response = self.client.get(
                        f"http://{self.broker_addr}/consume",
                        params={
                            "stream": nameInfo["stream_name"],
                            "group": nameInfo["group_name"],
                            "consumer": user_consumer["username"]
                        },
                        headers={"Authorization": f"Bearer {self.token}"}
                    )

                    print("---------------" + str(response.status_code) + "------------")
                    
                    if response.status_code == 204:  # No Content
                        print(f"No new messages in stream {nameInfo['stream_name']}, moving to next stream...")
                        continue
                        # 如果没有新消息，不需要立即重试，直接继续下一个stream
                    
                    elif response.status_code == 200:
                        messages = response.json()
                        if not messages:
                            print(f"No new message now in stream {nameInfo['stream_name']}, moving to next stream...")
                            continue
                        else:
                            print(f"Consumed messages from {nameInfo['stream_name']}.")
                            for message in messages:
                                self.ack(message)
                        # 继续下一个stream
                    
                    else:
                        # 出现其他响应码，记录错误并继续尝试
                        print(f"Failed to receive messages from stream {nameInfo['stream_name']}, status code: {response.status_code}")
                        time.sleep(1)  # 等待一段时间后重试

                except Exception as e:
                    # 捕获异常并记录日志
                    print(f"Error during message consumption from stream {nameInfo['stream_name']}: {str(e)}")
                    time.sleep(1)  # 等待一段时间后重试

            # 在完成一次对所有streams的遍历后，可以选择等待一段时间再开始新一轮遍历
            # time.sleep(1)  # 等待一段时间后再进行下一轮遍历



class ConsumerUser(HttpUser):
    tasks = [ConsumerTasks]
    wait_time = between(1, 2)
    host = proxy_url
