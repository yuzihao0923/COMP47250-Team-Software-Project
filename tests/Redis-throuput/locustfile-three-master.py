from locust import User, TaskSet, task, between
import redis
import json
import time

# 配置 Redis 集群连接
REDIS_NODES = [
    "redis://localhost:6381", 
    "redis://localhost:6382",
    "redis://localhost:6383",
    "redis://localhost:6384",
    "redis://localhost:6385",
    "redis://localhost:6386",
]

class RedisClient:
    def __init__(self):
        self.redis_cluster = redis.RedisCluster(startup_nodes=[{"host": "localhost", "port": 6379}])

    def set(self, key, value):
        self.redis_cluster.set(key, value)

    def get(self, key):
        return self.redis_cluster.get(key)

class RedisBehavior(TaskSet):
    def on_start(self):
        # Initialize Redis client
        self.client = RedisClient()

    @task
    def set_value(self):
        key = "test_key"
        value = "test_value"
        self.client.set(key, value)

    @task
    def get_value(self):
        key = "test_key"
        self.client.get(key)

class WebsiteUser(User):
    tasks = [RedisBehavior]
    wait_time = between(1, 2)
