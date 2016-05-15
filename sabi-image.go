package main

import (
    "fmt"
    "html"
    "log"
    "net/http"

    "github.com/streadway/amqp"
)

func main() {
    connectToBroker();

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
    })

    log.Fatal(http.ListenAndServe(":8080", nil))
}

func failOnError(err error, msg string) {
    if err != nil {
        log.Fatalf("%s: %s", msg, err)
        panic(fmt.Sprintf("%s: %s", msg, err))
    }
}

func connectToBroker() {
    conn, err := amqp.Dial("amqp://guest:guest@127.0.0.1:5672/")
    failOnError(err, "Failed to connect to RabbitMQ")
    defer conn.Close()

    ch, err := conn.Channel()
    failOnError(err, "Failed to open a channel")
    defer ch.Close()

    q, err := ch.QueueDeclare(
        "routingkey", // name
        false,   // durable
        false,   // delete when usused
        false,   // exclusive
        false,   // no-wait
        nil,     // arguments
    )
    failOnError(err, "Failed to declare a queue")

    body := "test one"
    err = ch.Publish(
        "test_exchange", // exchange
        q.Name,          // routing key
        false,           // mandatory
        false,           // immediate
        amqp.Publishing{
            ContentType: "text/plain",
            Body:        []byte(body),
        })
    failOnError(err, "Failed to publish a message")
}
