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
        # pass
        while True:
            try:
                # 发送请求以消费消息
                nameInfoList = self.read_from_json_file()
                nameInfo = random.choice(nameInfoList)
                print(nameInfo)
                response = self.client.get(
                    f"http://{self.broker_addr}/consume",
                    params={
                        "stream": nameInfo["stream_name"],
                        "group": nameInfo["group_name"],
                        "consumer": user_consumer["username"]
                    },
                    headers={"Authorization": f"Bearer {self.token}"}
                )

                print("---------------"+str(response.status_code)+"------------")
                
                if response.status_code == 204:  # No Content
                    print("No new messages, retrying...")
                    time.sleep(1)  # Wait for a bit before retrying
                
                elif response.status_code == 200:
                    messages = response.json()
                    if messages is None:
                        print("No new message now, please wait.")
                    else:
                        print(f"Consumed messages.")
                        for message in messages:
                            self.ack(message)
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
