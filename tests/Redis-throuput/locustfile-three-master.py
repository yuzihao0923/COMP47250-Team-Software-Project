from locust import User, TaskSet, task, between, events
from rediscluster import RedisCluster
import random
import string
import time

# 配置 Redis 集群连接
REDIS_NODES = [
    {"host": "localhost", "port": "6381"},
    {"host": "localhost", "port": "6382"},
    {"host": "localhost", "port": "6383"},
    {"host": "localhost", "port": "6384"},
    {"host": "localhost", "port": "6385"},
    {"host": "localhost", "port": "6386"},
]

class RedisClient:
    def __init__(self, environment):
        self.redis_cluster = RedisCluster(startup_nodes=REDIS_NODES, decode_responses=True)
        self.environment = environment
        print("Connected to Redis cluster.")

    def set(self, key, value):
        start_time = time.perf_counter_ns()
        try:
            self.redis_cluster.set(key, value)
            elapsed_time = time.perf_counter_ns() - start_time  # 纳秒
            self.environment.events.request.fire(
                request_type="SET", name="set", response_time=elapsed_time, response_length=0, exception=None
            )
        except Exception as e:
            elapsed_time = time.perf_counter_ns() - start_time
            self.environment.events.request.fire(
                request_type="SET", name="set", response_time=elapsed_time, response_length=0, exception=e
            )

    def get(self, key):
        start_time = time.perf_counter_ns()
        try:
            value = self.redis_cluster.get(key)
            elapsed_time = time.perf_counter_ns() - start_time  # 纳秒
            self.environment.events.request.fire(
                request_type="GET", name="get", response_time=elapsed_time, response_length=len(value) if value else 0, exception=None
            )
        except Exception as e:
            elapsed_time = time.perf_counter_ns() - start_time
            self.environment.events.request.fire(
                request_type="GET", name="get", response_time=elapsed_time, response_length=0, exception=e
            )

class RedisBehavior(TaskSet):
    def on_start(self):
        # Initialize Redis client with environment for reporting
        self.redis_client = RedisClient(self.user.environment)

    @task
    def set_value(self):
        key = ''.join(random.choices(string.ascii_letters + string.digits, k=10))  # Generate random key
        value = ''.join(random.choices(string.ascii_letters + string.digits, k=100))  # Generate random value
        self.redis_client.set(key, value)

    @task
    def get_value(self):
        key = ''.join(random.choices(string.ascii_letters + string.digits, k=10))  # Generate random key
        self.redis_client.get(key)

class WebsiteUser(User):
    tasks = [RedisBehavior]
    wait_time = between(1, 2)
