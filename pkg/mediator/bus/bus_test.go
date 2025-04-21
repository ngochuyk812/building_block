package bus_core

import (
	"context"
	"testing"

	"github.com/ngochuyk812/building_block/infrastructure/helpers"
	auth_context "github.com/ngochuyk812/building_block/pkg/auth"
	"github.com/ngochuyk812/building_block/pkg/mediator"
)

type CreateSiteCommand struct {
	Name string
}
type CreateSiteCommandResponse struct {
	Name string
}

func TestMediator(t *testing.T) {
	mediator := mediator.NewMediator()
	RegisterHandler(mediator, CreateSiteCommand{}, &CreateSiteHandler{})
	ctx := context.Background()
	ctx = helpers.NewContext(ctx, helpers.AuthContextKey, &auth_context.AuthContext{
		IdSite:   "123",
		UserName: "ngochuy",
	})
	req := CreateSiteCommand{Name: "domain.com"}
	rs, err := Send[CreateSiteCommand, CreateSiteCommandResponse](mediator, ctx, req)
	if err != nil {
		t.Error(err)
	}
	if rs.Name != req.Name {
		t.Errorf("cannot valid name: %s", req.Name)
	}

}

type CreateSiteHandler struct {
}

var _ IHandler[CreateSiteCommand, CreateSiteCommandResponse] = (*CreateSiteHandler)(nil)

func (h *CreateSiteHandler) Handle(ctx context.Context, cmd CreateSiteCommand) (CreateSiteCommandResponse, error) {

	return CreateSiteCommandResponse{
		Name: cmd.Name,
	}, nil
}
