package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"

	"time"
)

// CaseyTrigger holds the schema definition for the CaseyTrigger entity.
type CaseyTrigger struct {
	ent.Schema
}

// Fields of the CaseyTrigger.
func (CaseyTrigger) Fields() []ent.Field {

	return []ent.Field{

		field.Int64("id").SchemaType(map[string]string{
			dialect.MySQL: "bigint", // Override MySQL.
		}).Unique(),

		field.String("executor").SchemaType(map[string]string{
			dialect.MySQL: "varchar(255)", // Override MySQL.
		}),

		field.String("hostname").SchemaType(map[string]string{
			dialect.MySQL: "varchar(255)", // Override MySQL.
		}),

		field.Text("command").SchemaType(map[string]string{
			dialect.MySQL: "text", // Override MySQL.
		}),

		field.Text("response").SchemaType(map[string]string{
			dialect.MySQL: "text", // Override MySQL.
		}).Optional(),

		field.Time("executed_at").SchemaType(map[string]string{
			dialect.MySQL: "timestamp", // Override MySQL.
		}).Optional().Default(time.Now),
	}

}

// Edges of the CaseyTrigger.
func (CaseyTrigger) Edges() []ent.Edge {
	return nil
}
func (CaseyTrigger) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "casey_trigger"},
	}
}
