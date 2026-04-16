package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// RedeemCodeClaim records a user's successful claim of a welfare redeem code.
type RedeemCodeClaim struct {
	ent.Schema
}

func (RedeemCodeClaim) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "redeem_code_claims"},
	}
}

func (RedeemCodeClaim) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("redeem_code_id"),
		field.Int64("user_id"),
		field.Float("amount").
			SchemaType(map[string]string{dialect.Postgres: "decimal(20,8)"}),
		field.Time("claimed_at").
			Default(time.Now).
			SchemaType(map[string]string{dialect.Postgres: "timestamptz"}),
	}
}

func (RedeemCodeClaim) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("redeem_code", RedeemCode.Type).
			Ref("claims").
			Field("redeem_code_id").
			Unique().
			Required(),
		edge.From("user", User.Type).
			Ref("redeem_code_claims").
			Field("user_id").
			Unique().
			Required(),
	}
}

func (RedeemCodeClaim) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("redeem_code_id", "user_id").Unique(),
		index.Fields("user_id", "claimed_at"),
		index.Fields("redeem_code_id"),
	}
}
