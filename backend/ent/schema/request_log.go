// Package schema 定义 Ent ORM 的数据库 schema。
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

// RequestLog 定义请求日志实体的 schema。
//
// 记录每次 API 请求的元数据，用于会话可视化（不包含正文）。
type RequestLog struct {
	ent.Schema
}

// Annotations 返回 schema 的注解配置。
func (RequestLog) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "request_logs"},
	}
}

// Fields 定义请求日志实体的所有字段。
func (RequestLog) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("user_id"),
		field.Int64("api_key_id"),
		field.String("request_id").
			MaxLen(64).
			Optional().
			Nillable(),
		field.String("model").
			MaxLen(100).
			NotEmpty(),
		field.String("inbound_endpoint").
			MaxLen(128).
			Optional().
			Nillable(),
		field.String("upstream_endpoint").
			MaxLen(128).
			Optional().
			Nillable(),
		field.String("method").
			MaxLen(8).
			Optional().
			Nillable(),
		field.Int("status_code").
			Optional().
			Nillable(),
		field.String("error_code").
			MaxLen(64).
			Optional().
			Nillable(),
		field.String("error_message").
			Optional().
			Nillable().
			SchemaType(map[string]string{dialect.Postgres: "text"}),

		// Token 与成本
		field.Int("input_tokens").
			Default(0),
		field.Int("output_tokens").
			Default(0),
		field.Float("total_cost").
			Default(0).
			SchemaType(map[string]string{dialect.Postgres: "decimal(20,10)"}),

		// 其他元数据
		field.Bool("stream").
			Default(false),
		field.Int("duration_ms").
			Optional().
			Nillable(),
		field.Int("first_token_ms").
			Optional().
			Nillable(),
		field.String("user_agent").
			MaxLen(512).
			Optional().
			Nillable(),
		field.String("ip_address").
			MaxLen(45).
			Optional().
			Nillable(),

		// 时间戳（只有 created_at，日志不可修改）
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			SchemaType(map[string]string{dialect.Postgres: "timestamptz"}),
	}
}

// Edges 定义请求日志实体的关联关系。
func (RequestLog) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("request_logs").
			Field("user_id").
			Required().
			Unique(),
		edge.From("api_key", APIKey.Type).
			Ref("request_logs").
			Field("api_key_id").
			Required().
			Unique(),
	}
}

// Indexes 定义数据库索引，优化查询性能。
func (RequestLog) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id"),
		index.Fields("api_key_id"),
		index.Fields("request_id"),
		index.Fields("model"),
		index.Fields("status_code"),
		index.Fields("created_at"),
		index.Fields("user_id", "created_at"),
		index.Fields("api_key_id", "created_at"),
	}
}
