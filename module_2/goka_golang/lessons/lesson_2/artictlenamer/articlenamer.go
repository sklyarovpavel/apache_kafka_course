package artictlenamer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/lovoo/goka"
	"github.com/lovoo/goka/codec"
)

var (
	Group goka.Group = "namer"
)

type ValueCodec struct {
	codec.String
}

func replaceName(ctx goka.Context, msg interface{}) {
	ctx.SetValue(msg.(string))
}

type UserLikeArticle struct {
	Like    bool   //лайк поставленный пользователем
	UserId  int    // id пользователя
	Article string // Article которой был поставлен лайк
}

// UserLikeArticleCodec позволяет сериализовать и десериализовать пользователя в/из групповой таблицы.
type UserLikeArticleCodec struct{}

// Encode переводит user в []byte
func (uc *UserLikeArticleCodec) Encode(value any) ([]byte, error) {
	if _, isUser := value.(*UserLikeArticle); !isUser {
		return nil, fmt.Errorf("Тип должен быть *UserLike, получен %T", value)
	}
	return json.Marshal(value)
}

// Decode переводит user из []byte в структуру user.
func (uc *UserLikeArticleCodec) Decode(data []byte) (any, error) {
	var (
		c   UserLikeArticle
		err error
	)
	err = json.Unmarshal(data, &c)
	if err != nil {
		return nil, fmt.Errorf("Ошибка десериализации: %v", err)
	}
	return &c, nil
}

func RunArticleNamer(broker []string, inputStream goka.Stream) {
	g := goka.DefineGroup(Group,
		goka.Input(inputStream, new(ValueCodec), replaceName),
		goka.Persist(new(ValueCodec)),
	)
	p, err := goka.NewProcessor(broker, g)
	if err != nil {
		log.Fatal(err)
	}
	err = p.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
