package main

import (
	"github.com/AlexBlackNn/kafka-stream-go/stream/lesson_2/artictlenamer"
	"github.com/AlexBlackNn/kafka-stream-go/stream/lesson_2/blocker"
	"github.com/AlexBlackNn/kafka-stream-go/stream/lesson_2/emitter"
	"github.com/AlexBlackNn/kafka-stream-go/stream/lesson_2/filter"
	"github.com/AlexBlackNn/kafka-stream-go/stream/lesson_2/user"
	"github.com/lovoo/goka"
)

var brokers = []string{"127.0.0.1:29092"}

var (
	Emitter2FilterTopic       goka.Stream = "emitter2filter-stream"
	Filter2UserProcessorTopic goka.Stream = "filter2userprocessor-stream"
	BlockerTopic              goka.Stream = "blocker-stream"
	ArticleNamerTopic         goka.Stream = "namer-stream"
)

func main() {
	go blocker.RunBlocker(brokers, BlockerTopic)
	go artictlenamer.RunArticleNamer(brokers, ArticleNamerTopic)
	go emitter.RunEmitter(brokers, Emitter2FilterTopic)
	go filter.RunFilter(brokers, Emitter2FilterTopic, Filter2UserProcessorTopic)
	go user.RunUserProcessor(brokers, Filter2UserProcessorTopic)
	user.RunUserView(brokers)
}
