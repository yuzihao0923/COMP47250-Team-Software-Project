import redis
import time
from datetime import datetime

# Redis cluster configuration
cluster_nodes = [
    {"host": "127.0.0.1", "port": 6381},
    {"host": "127.0.0.1", "port": 6382},
    {"host": "127.0.0.1", "port": 6383},
    {"host": "127.0.0.1", "port": 6384},
    {"host": "127.0.0.1", "port": 6385},
    {"host": "127.0.0.1", "port": 6386}
]

def get_master_info(r):
    try:
        info = r.info('replication')
        return {
            'role': info['role'],
            'master_host': info.get('master_host'),
            'master_port': info.get('master_port'),
            'master_replid': info.get('master_replid'),
            'master_repl_offset': info.get('master_repl_offset')
        }
    except Exception as e:
        print(f"Error fetching replication info: {e}")
        return None

def monitor_cluster(nodes):
    previous_master_info = None
    while True:
        for node in nodes:
            try:
                r = redis.Redis(host=node['host'], port=node['port'], decode_responses=True)
                master_info = get_master_info(r)

                if master_info:
                    if previous_master_info:
                        if master_info['master_replid'] != previous_master_info['master_replid']:
                            now = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
                            print(f"{now} Master changed from {previous_master_info['master_replid']} to {master_info['master_replid']}")
                            # Calculate time taken for master switch
                            if previous_master_info['master_replid'] != master_info['master_replid']:
                                print(f"Master changed from {previous_master_info['master_replid']} to {master_info['master_replid']}. Time taken: {time.time() - previous_time:.2f} seconds")
                                previous_time = time.time()
                    previous_master_info = master_info
                    previous_time = time.time()  # Update previous_time for the next change

            except Exception as e:
                print(f"Error connecting to node {node['host']}:{node['port']}: {e}")

        time.sleep(60)  # Check every minute

monitor_cluster(cluster_nodes)
