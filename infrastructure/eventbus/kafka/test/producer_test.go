package mockup

import (
	"context"
	"strings"
	"testing"

	"github.com/ngochuyk812/building_block/infrastructure/eventbus/kafka"
	"github.com/ngochuyk812/building_block/infrastructure/helpers"
	auth_context "github.com/ngochuyk812/building_block/pkg/auth"
)

func TestNewProducer(t *testing.T) {
	producer, err := kafka.NewProceduer(strings.Split("localhost:9092", ","), "topic-test")
	if err != nil {
		t.Error(err)
	}
	if producer == nil {
		t.Errorf("producer not null")
	}
	ctx := context.Background()
	ctx = helpers.NewContext(ctx, helpers.AuthContextKey, &auth_context.AuthContext{
		IdSite:     "12323",
		IdAuthUser: "scscsc",
		Roles:      []string{"sdsd", "sdsd"},
		UserAgent:  "xccxcxcxcxcxcxc",
		UserIP:     "231312312.132.123.12.312",
		UserName:   "Ngochuy",
		Email:      "email@gmail.cokm",
	})
	err = producer.Publish(ctx, &EventSendMail{
		Content: "Content aa",
		Title:   "TItle a",
	})
	if err != nil {
		t.Error(err)
	}
}
