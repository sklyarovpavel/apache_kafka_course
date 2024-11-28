package user

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lovoo/goka"
)

var (
	Group goka.Group = "user-like"
)

// Like — объект, который хранится в групповой таблице процессора.
type Like struct {
	Like   bool //лайк поставленный пользователем
	UserId int  // id пользователя
	PostId int  // id статьи которой был поставлен лайк
}

// LikeCodec позволяет сериализовать и десериализовать пользователя в/из групповой таблицы.
type LikeCodec struct{}

// Encode переводит user в []byte
func (uc *LikeCodec) Encode(value any) ([]byte, error) {
	if _, isUser := value.(*Like); !isUser {
		return nil, fmt.Errorf("Тип должен быть *UserLike, получен %T", value)
	}
	return json.Marshal(value)
}

// Decode переводит user из []byte в структуру user.
func (uc *LikeCodec) Decode(data []byte) (any, error) {
	var (
		c   Like
		err error
	)
	err = json.Unmarshal(data, &c)
	if err != nil {
		return nil, fmt.Errorf("Ошибка десериализации: %v", err)
	}
	return &c, nil
}

type Post struct {
	PostLike map[string]bool
}

// PostCodec позволяет сериализовать и десериализовать пользователя в/из групповой таблицы.
type PostCodec struct{}

// Encode переводит user в []byte
func (up *PostCodec) Encode(value any) ([]byte, error) {
	if _, isUserPost := value.(*Post); !isUserPost {
		return nil, fmt.Errorf("Тип должен быть *user, получен %T", value)
	}
	return json.Marshal(value)
}

// Decode переводит user из []byte в структуру user.
func (up *PostCodec) Decode(data []byte) (any, error) {
	var (
		p   Post
		err error
	)
	err = json.Unmarshal(data, &p)
	if err != nil {
		return nil, fmt.Errorf("ошибка десериализации: %v", err)
	}
	return &p, nil
}

func process(ctx goka.Context, msg any) {
	var userLike *artictlenamer.UserLikeArticle
	var ok bool
	var userPost *Post

	if userLike, ok = msg.(*artictlenamer.UserLikeArticle); !ok || userLike == nil {
		return
	}
	fmt.Println("111111112", userLike)
	if val := ctx.Value(); val != nil {
		userPost = val.(*Post)
	} else {
		userPost = &Post{PostLike: make(map[string]bool)}
	}

	userPost.PostLike[userLike.Article] = userLike.Like

	ctx.SetValue(userPost)
	log.Printf("[proc] key: %s,  msg: %v, data in group_table %v \n", ctx.Key(), userLike, userPost)
}

func RunUserProcessor(brokers []string, inputStream goka.Stream) {
	g := goka.DefineGroup(Group,
		goka.Input(inputStream, new(artictlenamer.UserLikeArticleCodec), process),
		goka.Persist(new(PostCodec)),
	)
	p, err := goka.NewProcessor(brokers, g)
	if err != nil {
		log.Fatal(err)
	}
	err = p.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

func RunUserView(brokers []string) {
	view, err := goka.NewView(brokers,
		goka.GroupTable(Group),
		new(PostCodec),
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
