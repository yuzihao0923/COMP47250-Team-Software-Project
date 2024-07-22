from locust import HttpUser, TaskSet, task, between
import requests
import json
import random
import time
import base64

proxy_url = "http://localhost:8888"
user = {"username": "p1", "password": "123"}
MaxRetryCount = 3
RetryInterval = 2

class ProducerTasks(TaskSet):
    def on_start(self):
        self.broker_addr = self.get_broker_address()
        self.token = self.authenticate_user()

    def get_broker_address(self):
        response = requests.get(f"{proxy_url}/get-broker")
        response.raise_for_status()
        address = response.json()["address"]
        print(f"Broker address: {address}")
        return response.json()["address"]

    def authenticate_user(self):
        broker_login_url = f"http://{self.broker_addr}/login"
        response = self.client.post(
            broker_login_url,
            json={"username": user["username"], "password": user["password"]}
        )
        response.raise_for_status()
        data = response.json()
        if data["role"] == "producer":
            return data["token"]
        else:
            raise Exception(f"User {user['username']} is not a producer")

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
        for i in range(1):
            payload = f"Hello {i}".encode()
            str_payload = base64.b64encode(payload).decode('utf-8')
            self.send_message(self.broker_addr, "mystream", str_payload, token)

class ProducerUser(HttpUser):
    tasks = [ProducerTasks]
    wait_time = between(1, 2)
    host = proxy_url

if __name__ == "__main__":
    import os
    os.system("locust -f producerTest.py")
