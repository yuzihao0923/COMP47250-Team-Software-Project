from locust import HttpUser, TaskSet, task, between
import requests
import json
import time
import random
import string

proxy_url = "http://localhost:8888"
user_consumer = {"username": "c1", "password": "123"}

class ConsumerTasks(TaskSet):
    def on_start(self):
        self.broker_addr = self.get_broker_address()
        self.token = self.authenticate_user(user_consumer)
        self.stream_name = ''.join(random.choices(string.ascii_letters, k=10))
        self.group_name = ''.join(random.choices(string.ascii_letters, k=10))
        self.write_to_json_file()
        self.register()
    
    def read_from_json_file(self):
        try:
            with open("messages.json", "r") as file:
                data = json.load(file)
                return data
        except (IOError, json.JSONDecodeError) as e:
            print(f"Error reading JSON file: {e}")
            return []

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
    
    def write_to_json_file(self):
        message_data = {
            "stream_name": self.stream_name,
            "group_name": self.group_name
        }
        try:
            with open("messages.json", "r") as file:
                data = json.load(file)
        except (IOError, json.JSONDecodeError):
            data = []
        
        data.append(message_data)
        
        with open("messages.json", "w") as file:
            json.dump(data, file, indent=4)

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
    @task
    def consume_messages(self):
        while True:
            try:
                # 发送请求以消费消息
                nameInfo=self.read_from_json_file()
                response = self.client.get(
                    f"http://{self.broker_addr}/consume",
                    params={
                        "stream": nameInfo[0],
                        "group": nameInfo[1],
                        "consumer": user_consumer["username"]
                    },
                    headers={"Authorization": f"Bearer {self.token}"}
                )
                
                if response.status_code == 204:  # No Content
                    print("No new messages, retrying...")
                    time.sleep(1)  # Wait for a bit before retrying
                
                elif response.status_code == 200:
                    messages = response.json()
                    if len(messages) > 0:
                        print(f"Consumed {len(messages)} messages.")
                    else:
                        print("No new message now, please wait.")
                    break
                
                else:
                    print(f"Failed to receive messages, status code: {response.status_code}")
                    time.sleep(1)  # Wait for a bit before retrying

            except Exception as e:
                print(f"Error during message consumption: {str(e)}")
                time.sleep(1)  # Wait for a bit before retrying


class ConsumerUser(HttpUser):
    tasks = [ConsumerTasks]
    wait_time = between(1, 2)
    host = proxy_url
