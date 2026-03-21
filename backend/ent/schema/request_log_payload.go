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

// RequestLogPayload 定义请求日志正文快照实体。
//
// 将大字段从 request_logs 主表拆出，避免列表查询变重。
type RequestLogPayload struct {
	ent.Schema
}

func (RequestLogPayload) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "request_log_payloads"},
	}
}

func (RequestLogPayload) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("request_log_id"),
		field.Bytes("request_body").
			Optional(),
		field.String("request_body_encoding").
			MaxLen(16).
			Optional().
			Nillable(),
		field.Int("request_body_bytes").
			Optional().
			Nillable(),
		field.Bool("request_body_truncated").
			Default(false),
		field.Bytes("response_body").
			Optional(),
		field.String("response_body_encoding").
			MaxLen(16).
			Optional().
			Nillable(),
		field.Int("response_body_bytes").
			Optional().
			Nillable(),
		field.Bool("response_body_truncated").
			Default(false),
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			SchemaType(map[string]string{dialect.Postgres: "timestamptz"}),
	}
}

func (RequestLogPayload) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("request_log", RequestLog.Type).
			Ref("payload").
			Field("request_log_id").
			Required().
			Unique(),
	}
}

func (RequestLogPayload) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("request_log_id").
			Unique(),
		index.Fields("created_at"),
	}
}
