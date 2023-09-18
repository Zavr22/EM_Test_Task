package main

import "EMTestTask/pkg/kafka"

func main() {

	go kafka.ListenToKafkaTopic()

	select {}
}
