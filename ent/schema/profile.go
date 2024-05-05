package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Profile struct {
	ent.Schema
}

func (Profile) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

func (Profile) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			StorageKey("oid"),
		field.String("user_id").Unique(),
		field.String("full_name").Default(""),
		field.String("email"),
		field.Time("email_confirmed_at").
			Default(nil),
		field.String("phone"),
		field.Time("phone_confirmed_at").
			Default(nil),
		field.String("language").Default("en"),
		field.Bool("is_sso_user").Default(false),
	}
}

func (Profile) Edges() []ent.Edge {
	return nil
}
