package directive

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/DaisukeMatsumoto0925/backend2/graph/generated"
	gmodel "github.com/DaisukeMatsumoto0925/backend2/graph/model"
	"github.com/jinzhu/gorm"
)

type Directive struct {
	db *gorm.DB
}

func New(db *gorm.DB) generated.DirectiveRoot {
	d := &Directive{db}
	return generated.DirectiveRoot{
		HasRole: d.HasRole,
	}
}

func (d *Directive) HasRole(ctx context.Context, obj interface{}, next graphql.Resolver, role gmodel.Role) (interface{}, error) {
	token := getToken(ctx)
	if token == nil {
		return nil, fmt.Errorf("ACCESS DENIED")
	}
	if *token != role.String() {
		return nil, fmt.Errorf("ACCESS DENIED")
	}
	fmt.Println("authenticate here !")
	return next(ctx)
}

func getToken(ctx context.Context) *string {
	token := ctx.Value("token")
	// token := ctx.Value(tokenKey)
	if token, ok := token.(string); ok {
		return &token
	}
	return nil
}
