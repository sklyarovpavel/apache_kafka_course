package filter

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/AlexBlackNn/kafka-stream-go/stream/lesson_2/artictlenamer"
	"github.com/AlexBlackNn/kafka-stream-go/stream/lesson_2/blocker"
	"github.com/AlexBlackNn/kafka-stream-go/stream/lesson_2/user"
	"github.com/lovoo/goka"
)

var (
	filterGroup goka.Group = "filter"
)

func shouldDrop(ctx goka.Context) bool {
	v := ctx.Join(goka.GroupTable(blocker.Group))
	return v != nil && v.(*blocker.BlockValue).Blocked
}

func rename(ctx goka.Context, m *user.Like) *artictlenamer.UserLikeArticle {
	art := "no-name"
	var tw any
	if tw = ctx.Lookup(goka.GroupTable(artictlenamer.Group), strconv.Itoa(m.PostId)); tw != nil {
		art = tw.(string)
	}
	fmt.Println(&artictlenamer.UserLikeArticle{Like: m.Like, UserId: m.UserId, Article: art})
	return &artictlenamer.UserLikeArticle{Like: m.Like, UserId: m.UserId, Article: art}
}

func RunFilter(brokers []string, inputTopic goka.Stream, outputTopic goka.Stream) {
	g := goka.DefineGroup(filterGroup,
		goka.Input(inputTopic, new(user.LikeCodec), func(ctx goka.Context, msg interface{}) {
			if shouldDrop(ctx) {
				return
			}
			m := rename(ctx, msg.(*user.Like))
			fmt.Println("rename result", m)
			ctx.Emit(outputTopic, ctx.Key(), m)
		}),
		goka.Output(outputTopic, new(artictlenamer.UserLikeArticleCodec)),
		goka.Join(goka.GroupTable(blocker.Group), new(blocker.BlockValueCodec)),
		goka.Lookup(goka.GroupTable(artictlenamer.Group), new(artictlenamer.ValueCodec)),
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
