package main

import (
    "context"
    "encoding/json"
    "flag"
    "fmt"
    "log"
    "math/rand"
    "time"

    "github.com/segmentio/kafka-go"
)

func makeMetricsMessage(count int) map[string]interface{} {
    return map[string]interface{}{
        "host":      fmt.Sprintf("server-%d", rand.Intn(5)+1),
        "cpu":       rand.Float64() * 100,
        "ram":       rand.Intn(16000),
        "load":      rand.Float64() * 4,
        "count":     count,
        "timestamp": time.Now().UnixMilli(),
    }
}

func makeIoTMessage(count int) map[string]interface{} {
    return map[string]interface{}{
        "device":    fmt.Sprintf("sensor-%d", rand.Intn(50)),
        "temp":      rand.Float64()*15 + 10,
        "humidity":  rand.Intn(80),
        "battery":   rand.Intn(100),
        "count":     count,
        "timestamp": time.Now().UnixMilli(),
    }
}

func makeLogMessage(count int) map[string]interface{} {
    levels := []string{"INFO", "WARN", "ERROR", "DEBUG"}
    services := []string{"auth", "payments", "orders", "search"}

    return map[string]interface{}{
        "level":     levels[rand.Intn(len(levels))],
        "service":   services[rand.Intn(len(services))],
        "message":   "simulated log message",
        "count":     count,
        "timestamp": time.Now().UnixMilli(),
    }
}

func makePaymentMessage(count int) map[string]interface{} {
    statuses := []string{"AUTHORIZED", "FAILED", "PENDING"}

    return map[string]interface{}{
        "user":      rand.Intn(10000),
        "amount":    rand.Float64() * 200,
        "currency":  "USD",
        "status":    statuses[rand.Intn(len(statuses))],
        "count":     count,
        "timestamp": time.Now().UnixMilli(),
    }
}

func makeUserMessage(count int) map[string]interface{} {
    actions := []string{"login", "logout", "purchase", "view"}

    return map[string]interface{}{
        "user_id":   rand.Intn(10000),
        "action":    actions[rand.Intn(len(actions))],
        "device":    "mobile",
        "count":     count,
        "timestamp": time.Now().UnixMilli(),
    }
}

func buildMessage(mode string, count int) []byte {
    var msg map[string]interface{}

    switch mode {
    case "metrics":
        msg = makeMetricsMessage(count)
    case "iot":
        msg = makeIoTMessage(count)
    case "logs":
        msg = makeLogMessage(count)
    case "payments":
        msg = makePaymentMessage(count)
    case "users":
        msg = makeUserMessage(count)
    default:
        msg = map[string]interface{}{
            "count":     count,
            "timestamp": time.Now().UnixMilli(),
        }
    }

    data, _ := json.Marshal(msg)
    return data
}

func main() {
    rand.Seed(time.Now().UnixNano())

    broker := flag.String("broker", "localhost:9092", "Kafka broker address")
    topic := flag.String("topic", "test", "Kafka topic")
    interval := flag.Int("interval", 500, "Interval in ms")
    numPartitions := flag.Int("num-partitions", 1, "Number of partitions")
    mode := flag.String("mode", "metrics", "Message mode: metrics | iot | logs | payments | users")

    flag.Parse()

    writer := kafka.NewWriter(kafka.WriterConfig{
        Brokers:  []string{*broker},
        Topic:    *topic,
        Balancer: &kafka.Hash{},
    })

    ctx := context.Background()
    count := 0

    fmt.Println("Producer started; mode =", *mode)

    for {
        msg := kafka.Message{
            Partition: count % *numPartitions,
            Value:     buildMessage(*mode, count),
        }

        err := writer.WriteMessages(ctx, msg)
        if err != nil {
            log.Println("Kafka write error:", err)
        }

        count++
        time.Sleep(time.Duration(*interval) * time.Millisecond)
    }
}
