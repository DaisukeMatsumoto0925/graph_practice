package directive

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/DaisukeMatsumoto0925/backend2/graph/generated"
	gmodel "github.com/DaisukeMatsumoto0925/backend2/graph/model"
	"github.com/DaisukeMatsumoto0925/backend2/src/util/appcontext"
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
	token := appcontext.GetToken(ctx)
	if token == nil {
		return nil, fmt.Errorf("ACCESS DENIED")
	}
	if *token != role.String() {
		return nil, fmt.Errorf("ACCESS DENIED")
	}
	fmt.Println("authenticate here !")
	return next(ctx)
}
