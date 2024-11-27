package main

import (
	"flag"
	"log"

	"github.com/AlexBlackNn/kafka-stream-go/stream/lesson_2/blocker"
	"github.com/lovoo/goka"
)

var (
	user    = flag.String("user", "", "user to block")
	unblock = flag.Bool("unblock", false, "unblock user instead of blocking")
	broker  = flag.String("broker", "localhost:29092", "boostrap Kafka broker")
	stream  = flag.String("stream", "", "stream name")
)

func main() {
	flag.Parse()
	if *user == "" {
		log.Fatal("невозможно заблокировать пользователя ''")
	}
	emitter, err := goka.NewEmitter([]string{*broker}, goka.Stream(*stream), new(blocker.BlockEventCodec))
	if err != nil {
		log.Fatal(err)
	}
	defer emitter.Finish()

	err = emitter.EmitSync(*user, &blocker.BlockEvent{Unblock: *unblock})
	if err != nil {
		log.Fatal(err)
	}
}
