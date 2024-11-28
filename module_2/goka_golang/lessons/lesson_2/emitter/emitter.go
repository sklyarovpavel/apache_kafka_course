package emitter

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/AlexBlackNn/kafka-stream-go/stream/lesson_2/user"
	"github.com/lovoo/goka"
)

func RunEmitter(brokers []string, outputTopic goka.Stream) {
	// используется userLikeCodec так как отправляем структуру  UserLike
	emitter, err := goka.NewEmitter(brokers, outputTopic, new(user.LikeCodec))
	if err != nil {
		log.Fatal(err)
	}
	defer emitter.Finish()

	t := time.NewTicker(3 * time.Second)
	defer t.Stop()

	for range t.C {
		userId := rand.Intn(3)

		fakeUserLike := &user.Like{
			Like:   rand.Intn(2) == 1, // Случайное значение для лайка (true или false)
			UserId: userId,            // Случайный ID пользователя от 1 до 1000
			PostId: rand.Intn(5),      // Случайный ID статьи от 1 до 1000
		}

		err = emitter.EmitSync(fmt.Sprintf("user-%d", userId), fakeUserLike)
		if err != nil {
			log.Fatal(err)
		}
	}
}
