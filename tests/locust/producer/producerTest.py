from locust import HttpUser, TaskSet, task, between
import requests
import json
import time
import base64
import os
import string
import random

proxy_url = "http://localhost:8888"
user_producer = {"username": "p1", "password": "123"}
user_consumer = {"username": "c1", "password": "123"}
MaxRetryCount = 3
RetryInterval = 2

# Function to generate random string
def generate_random_string(n):
    return ''.join(random.choices(string.ascii_letters, k=n))

# Function to generate messages and write to file
def generate_and_write_messages(num_messages, file_path):
    messages = []
    for _ in range(num_messages):
        message = {
            "consumer_info": {
                "stream_name": generate_random_string(10),
                "group_name": generate_random_string(10),
            }
        }
        messages.append(message)
    
    with open(file_path, 'w') as file:
        json.dump(messages, file, indent=2)

class ProducerTasks(TaskSet):
    def on_start(self):
        self.broker_addr = self.get_broker_address()     
        self.token = self.authenticate_user(user_producer)           

    def get_broker_address(self):
        response = requests.get(f"{proxy_url}/get-broker")
        response.raise_for_status()
        address = response.json()["address"]
        print(f"Broker address: {address}")
        return address

    def authenticate_user(self, user):
        broker_login_url = f"http://{self.broker_addr}/login"
        response = self.client.post(
            broker_login_url,
            json={"username": user["username"], "password": user["password"]}
        )
        response.raise_for_status()
        data = response.json()
        if data["role"] == "producer" or data["role"] == "consumer":
            return data["token"]
        else:
            raise Exception(f"User {user['username']} is not authorized")

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
        for i in range(10):  # For demonstration, just sending 10 messages
            payload = f"Hello {i}".encode()
            str_payload = base64.b64encode(payload).decode('utf-8')
            self.send_message(self.broker_addr, "mystream", str_payload, token)

class ConsumerTasks(TaskSet):
    def on_start(self):
        self.broker_addr = self.get_broker_address()    
        self.token = self.authenticate_user(user_consumer)
        self.message_file_path = 'messages.json'
        self.all_registered = False

    def get_broker_address(self):
        response = requests.get(f"{proxy_url}/get-broker")
        response.raise_for_status()
        address = response.json()["address"]
        print(f"Broker address: {address}")
        return address

    def authenticate_user(self, user):
        broker_login_url = f"http://{self.broker_addr}/login"
        response = self.client.post(
            broker_login_url,
            json={"username": user["username"], "password": user["password"]}
        )
        response.raise_for_status()
        data = response.json()
        if data["role"] == "producer" or data["role"] == "consumer":
            return data["token"]
        else:
            raise Exception(f"User {user['username']} is not authorized")

    def register(self, broker_addr, stream_name, group_name, token):
        msg = {
            "type": "registration",
            "consumer_info": {
                "stream_name": stream_name,
                "group_name": group_name
            }
        }
        headers = {"Authorization": f"Bearer {self.token}", "Content-Type": "application/json"}
        response = self.client.post(f"http://{self.broker_addr}/register", json=msg, headers=headers)
        if response.status_code == 200:
            print(f"Consumer registered: {response.text}")
        else:
            print(f"Error registering consumer: {response.text}")

    def consume_message(self, broker_addr, stream_name, token):
        headers = {"Authorization": f"Bearer {token}"}
        for retry_count in range(MaxRetryCount):
            response = self.client.get(f"http://{broker_addr}/consume?stream_name={stream_name}", headers=headers)
            if response.status_code == 200:
                # Removed print statement to avoid logging the consumed message
                return response.json()['payload']
            print(f"Error consuming message (attempt {retry_count + 1}/{MaxRetryCount}): {response.text}")
            time.sleep(RetryInterval)
        print(f"Failed to consume message after {MaxRetryCount} attempts")

    @task
    def consume_messages(self):
        token = self.token

        if not self.all_registered:
            # Generate and write messages
            generate_and_write_messages(10, self.message_file_path)  # Adjust the number as needed
            
            # Load messages from file
            with open(self.message_file_path, 'r') as file:
                self.messages = json.load(file)
            
            # Register consumers (all at once)
            self.message_index = 0
            # Store the stream names and group names for registration
            registrations = []
            
            while self.message_index < len(self.messages):
                message = self.messages[self.message_index]
                stream_name = message['consumer_info']['stream_name']
                group_name = message['consumer_info']['group_name']
                registrations.append((stream_name, group_name))  # Save for later registration
                self.message_index += 1
            
            # Now register all consumers at once
            for stream_name, group_name in registrations:
                self.register(self.broker_addr, stream_name, group_name, token)
            
            self.all_registered = True
            print("All consumers registered")
        
        else:
            # Consume messages
            self.message_index = 0
            while self.message_index < len(self.messages):
                message = self.messages[self.message_index]
                stream_name = message['consumer_info']['stream_name']
                self.consume_message(self.broker_addr, stream_name, token)
                self.message_index += 1
            print("All messages consumed")



class ProducerUser(HttpUser):
    tasks = [ProducerTasks]
    wait_time = between(1, 2)
    host = proxy_url

class ConsumerUser(HttpUser):
    tasks = [ConsumerTasks]
    wait_time = between(1, 2)
    host = proxy_url

if __name__ == "__main__":
    os.system("locust -f producer_consumer_test.py")
