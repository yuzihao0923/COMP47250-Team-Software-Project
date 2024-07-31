from locust import HttpUser, TaskSet, task, between
import base64
import json
import time
import os
import requests
import random

proxy_url = "http://localhost:8888"
user_producer = {"username": "p1", "password": "123"}
MaxRetryCount = 3
RetryInterval = 2

class ProducerTasks(TaskSet):
    def on_start(self):
        self.broker_addr = self.get_broker_address()
        self.token = self.authenticate_user(user_producer)
    
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

    def read_from_json_file(self):
        try:
            with open("message.json", "r") as file:
                lines = file.readlines()
                if lines:
                    # 从所有行中随机选择一行
                    random_line = random.choice(lines).strip()
                    message_data = json.loads(random_line)
                    return message_data.get("stream_name"), message_data.get("group_name")
        except (IOError, json.JSONDecodeError) as e:
            print(f"Error reading JSON file: {e}")
        return None, None

    def send_message(self, broker_addr, stream_name, payload, token):
        msg = {
            "type": "produce",
            "consumer_info": {"stream_name": stream_name},
            "payload": payload
        }
        headers = {"Authorization": f"Bearer {token}"}
        for retry_count in range(MaxRetryCount):
            response = self.client.post(f"http://{broker_addr}/produce", json=msg, headers=headers)
            if response.status_code == 200:
                print(f"Producer sent message: {payload}")
                return
            print(f"Error sending message (attempt {retry_count + 1}/{MaxRetryCount}): {response.text}")
            time.sleep(RetryInterval)
        print(f"Failed to send message after {MaxRetryCount} attempts")

    @task
    def produce_message(self):
        token = self.token
        stream_name, _ = self.read_from_json_file()
        if stream_name:
            payload = f"Hello {stream_name}".encode()
            str_payload = base64.b64encode(payload).decode('utf-8')
            self.send_message(self.broker_addr, stream_name, str_payload, token)
        else:
            print("Stream name not found. Ensure consumers are registered and message.json is updated.")

class ProducerUser(HttpUser):
    tasks = [ProducerTasks]
    wait_time = between(1, 2)
    host = proxy_url
