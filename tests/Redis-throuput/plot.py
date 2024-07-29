import os
import matplotlib.pyplot as plt

# 读取文件数据
def read_throughput_data(file_path):
    with open(file_path, 'r') as file:
        return [int(line.strip()) for line in file.readlines()]

# 生成折线图
def generate_plot(configs, data):
    plt.figure(figsize=(10, 6))

    for config, throughput in zip(configs, data):
        plt.plot(range(1, len(throughput) + 1), throughput, marker='o', label=config)

    plt.xlabel('Test Iteration')
    plt.ylabel('Throughput')
    plt.title('Throughput Comparison Across Different Redis Clusters')
    plt.legend()
    plt.grid(True)
    plt.savefig('throughput_comparison.png')
    plt.show()

# 主函数
def main():
    configs = ["redis-3-master", "redis-4-master", "redis-5-master"]
    data = []

    for config in configs:
        file_path = f'results/{config}_throughput.txt'
        if os.path.exists(file_path):
            data.append(read_throughput_data(file_path))
        else:
            print(f'File not found: {file_path}')

    if data:
        generate_plot(configs, data)

if __name__ == '__main__':
    main()
