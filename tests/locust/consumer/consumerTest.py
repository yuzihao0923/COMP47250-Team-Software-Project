from locust import HttpUser, TaskSet, task, between
import base64
import time

proxy_url = "http://localhost:8888"
user = {"username": "c1", "password": "123"}
MaxRetryCount = 3
RetryInterval = 2

class ConsumerTasks(TaskSet):
    def on_start(self):
        self.broker_addr = self.get_broker_address()
        self.token = self.authenticate_user()

    def get_broker_address(self):
        response = self.client.get(f"{proxy_url}/get-broker")
        response.raise_for_status()
        address = response.json()["address"]
        print(f"Broker address: {address}")
        return address

    def authenticate_user(self):
        broker_login_url = f"http://{self.broker_addr}/login"
        response = self.client.post(
            broker_login_url,
            json={"username": user["username"], "password": user["password"]}
        )
        response.raise_for_status()
        data = response.json()
        if data["role"] == "consumer":
            return data["token"]
        else:
            raise Exception(f"User {user['username']} is not a consumer")

    def receive_message(self, broker_addr, stream_name, token):
        msg = {
            "type": "consume",
            "stream_name": stream_name
        }
        headers = {"Authorization": f"Bearer {token}"}
        for retry_count in range(MaxRetryCount):
            with self.client.post(f"http://{broker_addr}/consume", json=msg, headers=headers, catch_response=True) as response:
                if response.status_code == 200:
                    payload = response.json()["payload"]
                    decoded_msg = base64.b64decode(payload).decode('utf-8')
                    print(f"Consumer received message: {decoded_msg}")
                    response.success()
                    return
                else:
                    print(f"Error receiving message (attempt {retry_count + 1}/{MaxRetryCount}): {response.text}")
                    response.failure(f"Error {response.status_code}")
                    time.sleep(RetryInterval)
        print(f"Failed to receive message after {MaxRetryCount} attempts")

    @task
    def consume_message(self):
        self.receive_message(self.broker_addr, "mystream", self.token)

class ConsumerUser(HttpUser):
    tasks = [ConsumerTasks]
    wait_time = between(1, 2)
    host = proxy_url

if __name__ == "__main__":
    import os
    os.system("locust -f consumerTest.py")
