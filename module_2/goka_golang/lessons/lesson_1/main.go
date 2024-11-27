package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/lovoo/goka"
)

var (
	brokers             = []string{"127.0.0.1:29092"}
	topic   goka.Stream = "user-likes-stream"
	group   goka.Group  = "user-like-group"
)

// UserLike — объект, который хранится в групповой таблице процессора.
type UserLike struct {
	Like   bool //лайк поставленный пользователем
	UserId int  // id пользователя
	PostId int  // id статьи которой был поставлен лайк
}

// userLikeCodec позволяет сериализовать и десериализовать пользователя в/из групповой таблицы.
type userLikeCodec struct{}

// Encode переводит user в []byte
func (uc *userLikeCodec) Encode(value any) ([]byte, error) {
	if _, isUser := value.(*UserLike); !isUser {
		return nil, fmt.Errorf("Тип должен быть *UserLike, получен %T", value)
	}
	return json.Marshal(value)
}

// Decode переводит user из []byte в структуру user.
func (uc *userLikeCodec) Decode(data []byte) (any, error) {
	var (
		c   UserLike
		err error
	)
	err = json.Unmarshal(data, &c)
	if err != nil {
		return nil, fmt.Errorf("Ошибка десериализации: %v", err)
	}
	return &c, nil
}

type UserPost struct {
	PostLike map[int]bool
}

// userPostCodec позволяет сериализовать и десериализовать пользователя в/из групповой таблицы.
type userPostCodec struct{}

// Encode переводит user в []byte
func (up *userPostCodec) Encode(value any) ([]byte, error) {
	if _, isUserPost := value.(*UserPost); !isUserPost {
		return nil, fmt.Errorf("Тип должен быть *user, получен %T", value)
	}
	return json.Marshal(value)
}

// Decode переводит user из []byte в структуру user.
func (up *userPostCodec) Decode(data []byte) (any, error) {
	var (
		p   UserPost
		err error
	)
	err = json.Unmarshal(data, &p)
	if err != nil {
		return nil, fmt.Errorf("Ошибка десериализации: %v", err)
	}
	return &p, nil
}

func runEmitter() {
	// используется userLikeCodec так как отправляем структуру  UserLike
	emitter, err := goka.NewEmitter(brokers, topic, new(userLikeCodec))
	if err != nil {
		log.Fatal(err)
	}
	defer emitter.Finish()

	t := time.NewTicker(100 * time.Millisecond)
	defer t.Stop()

	var i int
	for range t.C {
		userId := rand.Intn(3)

		fakeUserLike := &UserLike{
			Like:   rand.Intn(2) == 1, // Случайное значение для лайка (true или false)
			UserId: userId,            // Случайный ID пользователя от 1 до 1000
			PostId: rand.Intn(5),      // Случайный ID статьи от 1 до 1000
		}

		err = emitter.EmitSync(fmt.Sprintf("user-%d", userId), fakeUserLike)
		if err != nil {
			log.Fatal(err)
		}
		i++
	}
}

func process(ctx goka.Context, msg any) {

	var userLike *UserLike
	var ok bool
	var userPost *UserPost

	if userLike, ok = msg.(*UserLike); !ok || userLike == nil {
		return
	}

	if val := ctx.Value(); val != nil {
		userPost = val.(*UserPost)
	} else {
		userPost = &UserPost{PostLike: make(map[int]bool)}
	}

	userPost.PostLike[userLike.PostId] = userLike.Like

	ctx.SetValue(userPost)
	log.Printf("[proc] key: %s,  msg: %v, data in group_table %v \n", ctx.Key(), userLike, userPost)
}

func runProcessor() {
	g := goka.DefineGroup(group,
		goka.Input(topic, new(userLikeCodec), process),
		goka.Persist(new(userPostCodec)),
	)
	p, err := goka.NewProcessor(brokers,
		g,
		goka.WithConsumerGroupBuilder(goka.DefaultConsumerGroupBuilder),
	)
	if err != nil {
		log.Fatal(err)
	}
	err = p.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

func runView() {

	view, err := goka.NewView(brokers,
		goka.GroupTable(group),
		new(userPostCodec),
	)
	if err != nil {
		log.Fatal(err)
	}

	root := mux.NewRouter()
	root.HandleFunc("/{key}", func(w http.ResponseWriter, r *http.Request) {
		value, _ := view.Get(mux.Vars(r)["key"])
		data, _ := json.Marshal(value)
		w.Write(data)
	})
	log.Println("View opened at http://localhost:9095/")
	go func() {
		err = http.ListenAndServe(":9095", root)
		if err != nil {
			log.Fatal(err)
		}
	}()
	err = view.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	go runEmitter()
	go runProcessor()
	runView()
}
