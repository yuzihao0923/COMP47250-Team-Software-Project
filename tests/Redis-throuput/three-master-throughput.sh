#!/bin/bash

NUM_TESTS=10
TEST_DURATION=60
CLUSTER_TYPE="redis-3-master"

# 创建结果和日志目录
mkdir -p results
mkdir -p logs

# Function to run single test iteration
function run_test_iteration {
    ITERATION=$1

    echo "Running test iteration $ITERATION for Redis cluster: $CLUSTER_TYPE"

    # 手动启动 producer 和 consumer 之前运行此脚本

    # 运行测试的特定时间
    sleep $TEST_DURATION

    # 收集和记录结果
    echo "Test iteration $ITERATION completed for Redis cluster: $CLUSTER_TYPE"
    echo "Collecting results..."

    # 收集日志或指标
    THROUGHPUT=$(grep -oP 'Processed message at \K.*' logs/broker.log | wc -l)
    echo "Iteration $ITERATION throughput: $THROUGHPUT"
    echo $THROUGHPUT >> results/${CLUSTER_TYPE}_throughput.txt

    # 手动停止 producer 和 consumer 之后运行此脚本
}

# Function to calculate average throughput
function calculate_average_throughput {
    TOTAL=0
    COUNT=0

    while read -r line; do
        TOTAL=$(($TOTAL + $line))
        COUNT=$(($COUNT + 1))
    done < results/${CLUSTER_TYPE}_throughput.txt

    AVERAGE=$(($TOTAL / $COUNT))
    echo "Average throughput for Redis cluster $CLUSTER_TYPE: $AVERAGE"
    echo $AVERAGE > results/${CLUSTER_TYPE}_average_throughput.txt
}

# 清除以前的结果
> results/${CLUSTER_TYPE}_throughput.txt

for ((i=1; i<=NUM_TESTS; i++)); do
    run_test_iteration $i
done

calculate_average_throughput
echo "Testing completed for Redis cluster: $CLUSTER_TYPE"
