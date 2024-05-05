package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

type TimeMixin struct {
	// We embed the `mixin.Schema` to avoid
	// implementing the rest of the methods.
	mixin.Schema
}

func (TimeMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").
			Immutable().
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

type Chart struct {
	ent.Schema
}

func (Chart) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

func (Chart) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			StorageKey("oid"),
		field.String("user_id").
			Default(""),
		field.Uint32("chart_id").
			Optional(),
		field.String("template_id").
			Default(""),
		field.String("client_id").
			Default(""),
		field.String("type").
			Default(""),
		field.String("name").
			Default(""),
		field.String("content").
			Default(""),
		field.String("symbol").
			Default(""),
		field.String("resolution").
			Default(""),
	}
}

func (Chart) Edges() []ent.Edge {
	return nil
}
