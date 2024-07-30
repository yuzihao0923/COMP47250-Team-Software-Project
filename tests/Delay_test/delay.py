import threading
import time
import redis
import json
from datetime import datetime
import requests
import xml.etree.ElementTree as ET

# 配置 Redis 集群的地址
redis_nodes = [
    {"host": "127.0.0.1", "port": 6381},
    {"host": "127.0.0.1", "port": 6382},
    {"host": "127.0.0.1", "port": 6383},
    {"host": "127.0.0.1", "port": 6384},
    {"host": "127.0.0.1", "port": 6385},
    {"host": "127.0.0.1", "port": 6386},
]

# 创建 Redis 集群客户端
redis_client = redis.RedisCluster(startup_nodes=redis_nodes, decode_responses=True)


# 生产者函数，通过HTTP请求向Broker发送消息
def producer(broker_url, stream_name, message, count):
    headers = {'Content-Type': 'application/json'}
    for _ in range(count):
        start_time = datetime.utcnow().isoformat()
        message["timestamp"] = start_time
        response = requests.post(f"{broker_url}/produce", headers=headers, data=json.dumps({
            "stream": stream_name,
            "message": message
        }))
        if response.status_code != 200:
            print(f"Failed to send message: {response.text}")
        time.sleep(0.01)  # 模拟发送间隔


# 消费者函数，从Redis流中读取消息
def consumer(stream_name, group_name, consumer_name, results, lock):
    while True:
        messages = redis_client.xreadgroup(group_name, consumer_name, {stream_name: '>'}, count=1, block=5000)
        for stream, msgs in messages:
            for msg_id, msg in msgs:
                end_time = datetime.utcnow().isoformat()
                start_time = msg[b'timestamp'].decode('utf-8')
                delay = (datetime.fromisoformat(end_time) - datetime.fromisoformat(start_time)).total_seconds()

                with lock:
                    results.append(delay)


# 创建消费者组
def create_consumer_group(stream_name, group_name):
    try:
        redis_client.xgroup_create(stream_name, group_name, mkstream=True)
    except redis.exceptions.ResponseError as e:
        if "BUSYGROUP Consumer Group name already exists" not in str(e):
            raise


# 运行测试
def run_test(broker_url, stream_name, group_name, producer_count, message_count, results):
    lock = threading.Lock()

    # 启动消费者线程
    consumer_thread = threading.Thread(target=consumer, args=(stream_name, group_name, "myconsumer", results, lock))
    consumer_thread.start()

    # 启动生产者线程
    producer_threads = []
    for i in range(producer_count):
        producer_thread = threading.Thread(target=producer,
                                           args=(broker_url, stream_name, {"producer": f"producer_{i}"}, message_count))
        producer_threads.append(producer_thread)
        producer_thread.start()

    # 等待所有生产者线程完成
    for thread in producer_threads:
        thread.join()

    # 等待一段时间，确保所有消息被消费
    time.sleep(10)

    # 关闭消费者线程（在真实场景中，您可能需要更优雅地处理线程终止）
    consumer_thread.join(timeout=1)


# 保存结果到XML文件
def save_results_to_xml(filename, results):
    root = ET.Element("Results")
    for result in results:
        entry = ET.SubElement(root, "Entry")
        entry.set("producers", str(result["producers"]))
        entry.set("messages", str(result["messages"]))
        entry.set("avg_delay", str(result["avg_delay"]))

    tree = ET.ElementTree(root)
    tree.write(filename, encoding="utf-8", xml_declaration=True)


# 主函数
def main():
    broker_url = "http://localhost:8080"  # Broker的地址
    stream_name = "mystream"
    group_name = "mygroup"

    # 创建消费者组
    create_consumer_group(stream_name, group_name)

    # 测试参数
    results = []

    for producer_count in range(1, 6):  # 从1到5倍增加并发量
        for message_count in range(100, 10001, 500):  # 从100到10000递增消息量
            print(f"Running test with {producer_count} producers and {message_count} messages")
            test_results = []
            run_test(broker_url, stream_name, group_name, producer_count, message_count, test_results)
            avg_delay = sum(test_results) / len(test_results) if test_results else 0
            results.append({
                "producers": producer_count,
                "messages": message_count,
                "avg_delay": avg_delay
            })

    # 保存结果到XML文件
    xml_filename = "test_results.xml"
    save_results_to_xml(xml_filename, results)
    print(f"Results saved to {xml_filename}")


if __name__ == "__main__":
    main()
