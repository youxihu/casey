// Code generated by ent, DO NOT EDIT.

package caseytrigger

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/youxihu/casey/internal/data/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int64) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int64) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int64) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int64) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int64) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int64) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int64) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int64) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int64) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldLTE(FieldID, id))
}

// Executor applies equality check predicate on the "executor" field. It's identical to ExecutorEQ.
func Executor(v string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldEQ(FieldExecutor, v))
}

// Hostname applies equality check predicate on the "hostname" field. It's identical to HostnameEQ.
func Hostname(v string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldEQ(FieldHostname, v))
}

// Command applies equality check predicate on the "command" field. It's identical to CommandEQ.
func Command(v string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldEQ(FieldCommand, v))
}

// Response applies equality check predicate on the "response" field. It's identical to ResponseEQ.
func Response(v string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldEQ(FieldResponse, v))
}

// ExecutedAt applies equality check predicate on the "executed_at" field. It's identical to ExecutedAtEQ.
func ExecutedAt(v time.Time) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldEQ(FieldExecutedAt, v))
}

// ExecutorEQ applies the EQ predicate on the "executor" field.
func ExecutorEQ(v string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldEQ(FieldExecutor, v))
}

// ExecutorNEQ applies the NEQ predicate on the "executor" field.
func ExecutorNEQ(v string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldNEQ(FieldExecutor, v))
}

// ExecutorIn applies the In predicate on the "executor" field.
func ExecutorIn(vs ...string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldIn(FieldExecutor, vs...))
}

// ExecutorNotIn applies the NotIn predicate on the "executor" field.
func ExecutorNotIn(vs ...string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldNotIn(FieldExecutor, vs...))
}

// ExecutorGT applies the GT predicate on the "executor" field.
func ExecutorGT(v string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldGT(FieldExecutor, v))
}

// ExecutorGTE applies the GTE predicate on the "executor" field.
func ExecutorGTE(v string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldGTE(FieldExecutor, v))
}

// ExecutorLT applies the LT predicate on the "executor" field.
func ExecutorLT(v string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldLT(FieldExecutor, v))
}

// ExecutorLTE applies the LTE predicate on the "executor" field.
func ExecutorLTE(v string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldLTE(FieldExecutor, v))
}

// ExecutorContains applies the Contains predicate on the "executor" field.
func ExecutorContains(v string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldContains(FieldExecutor, v))
}

// ExecutorHasPrefix applies the HasPrefix predicate on the "executor" field.
func ExecutorHasPrefix(v string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldHasPrefix(FieldExecutor, v))
}

// ExecutorHasSuffix applies the HasSuffix predicate on the "executor" field.
func ExecutorHasSuffix(v string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldHasSuffix(FieldExecutor, v))
}

// ExecutorEqualFold applies the EqualFold predicate on the "executor" field.
func ExecutorEqualFold(v string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldEqualFold(FieldExecutor, v))
}

// ExecutorContainsFold applies the ContainsFold predicate on the "executor" field.
func ExecutorContainsFold(v string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldContainsFold(FieldExecutor, v))
}

// HostnameEQ applies the EQ predicate on the "hostname" field.
func HostnameEQ(v string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldEQ(FieldHostname, v))
}

// HostnameNEQ applies the NEQ predicate on the "hostname" field.
func HostnameNEQ(v string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldNEQ(FieldHostname, v))
}

// HostnameIn applies the In predicate on the "hostname" field.
func HostnameIn(vs ...string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldIn(FieldHostname, vs...))
}

// HostnameNotIn applies the NotIn predicate on the "hostname" field.
func HostnameNotIn(vs ...string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldNotIn(FieldHostname, vs...))
}

// HostnameGT applies the GT predicate on the "hostname" field.
func HostnameGT(v string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldGT(FieldHostname, v))
}

// HostnameGTE applies the GTE predicate on the "hostname" field.
func HostnameGTE(v string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldGTE(FieldHostname, v))
}

// HostnameLT applies the LT predicate on the "hostname" field.
func HostnameLT(v string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldLT(FieldHostname, v))
}

// HostnameLTE applies the LTE predicate on the "hostname" field.
func HostnameLTE(v string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldLTE(FieldHostname, v))
}

// HostnameContains applies the Contains predicate on the "hostname" field.
func HostnameContains(v string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldContains(FieldHostname, v))
}

// HostnameHasPrefix applies the HasPrefix predicate on the "hostname" field.
func HostnameHasPrefix(v string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldHasPrefix(FieldHostname, v))
}

// HostnameHasSuffix applies the HasSuffix predicate on the "hostname" field.
func HostnameHasSuffix(v string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldHasSuffix(FieldHostname, v))
}

// HostnameEqualFold applies the EqualFold predicate on the "hostname" field.
func HostnameEqualFold(v string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldEqualFold(FieldHostname, v))
}

// HostnameContainsFold applies the ContainsFold predicate on the "hostname" field.
func HostnameContainsFold(v string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldContainsFold(FieldHostname, v))
}

// CommandEQ applies the EQ predicate on the "command" field.
func CommandEQ(v string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldEQ(FieldCommand, v))
}

// CommandNEQ applies the NEQ predicate on the "command" field.
func CommandNEQ(v string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldNEQ(FieldCommand, v))
}

// CommandIn applies the In predicate on the "command" field.
func CommandIn(vs ...string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldIn(FieldCommand, vs...))
}

// CommandNotIn applies the NotIn predicate on the "command" field.
func CommandNotIn(vs ...string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldNotIn(FieldCommand, vs...))
}

// CommandGT applies the GT predicate on the "command" field.
func CommandGT(v string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldGT(FieldCommand, v))
}

// CommandGTE applies the GTE predicate on the "command" field.
func CommandGTE(v string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldGTE(FieldCommand, v))
}

// CommandLT applies the LT predicate on the "command" field.
func CommandLT(v string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldLT(FieldCommand, v))
}

// CommandLTE applies the LTE predicate on the "command" field.
func CommandLTE(v string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldLTE(FieldCommand, v))
}

// CommandContains applies the Contains predicate on the "command" field.
func CommandContains(v string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldContains(FieldCommand, v))
}

// CommandHasPrefix applies the HasPrefix predicate on the "command" field.
func CommandHasPrefix(v string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldHasPrefix(FieldCommand, v))
}

// CommandHasSuffix applies the HasSuffix predicate on the "command" field.
func CommandHasSuffix(v string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldHasSuffix(FieldCommand, v))
}

// CommandEqualFold applies the EqualFold predicate on the "command" field.
func CommandEqualFold(v string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldEqualFold(FieldCommand, v))
}

// CommandContainsFold applies the ContainsFold predicate on the "command" field.
func CommandContainsFold(v string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldContainsFold(FieldCommand, v))
}

// ResponseEQ applies the EQ predicate on the "response" field.
func ResponseEQ(v string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldEQ(FieldResponse, v))
}

// ResponseNEQ applies the NEQ predicate on the "response" field.
func ResponseNEQ(v string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldNEQ(FieldResponse, v))
}

// ResponseIn applies the In predicate on the "response" field.
func ResponseIn(vs ...string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldIn(FieldResponse, vs...))
}

// ResponseNotIn applies the NotIn predicate on the "response" field.
func ResponseNotIn(vs ...string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldNotIn(FieldResponse, vs...))
}

// ResponseGT applies the GT predicate on the "response" field.
func ResponseGT(v string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldGT(FieldResponse, v))
}

// ResponseGTE applies the GTE predicate on the "response" field.
func ResponseGTE(v string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldGTE(FieldResponse, v))
}

// ResponseLT applies the LT predicate on the "response" field.
func ResponseLT(v string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldLT(FieldResponse, v))
}

// ResponseLTE applies the LTE predicate on the "response" field.
func ResponseLTE(v string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldLTE(FieldResponse, v))
}

// ResponseContains applies the Contains predicate on the "response" field.
func ResponseContains(v string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldContains(FieldResponse, v))
}

// ResponseHasPrefix applies the HasPrefix predicate on the "response" field.
func ResponseHasPrefix(v string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldHasPrefix(FieldResponse, v))
}

// ResponseHasSuffix applies the HasSuffix predicate on the "response" field.
func ResponseHasSuffix(v string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldHasSuffix(FieldResponse, v))
}

// ResponseIsNil applies the IsNil predicate on the "response" field.
func ResponseIsNil() predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldIsNull(FieldResponse))
}

// ResponseNotNil applies the NotNil predicate on the "response" field.
func ResponseNotNil() predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldNotNull(FieldResponse))
}

// ResponseEqualFold applies the EqualFold predicate on the "response" field.
func ResponseEqualFold(v string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldEqualFold(FieldResponse, v))
}

// ResponseContainsFold applies the ContainsFold predicate on the "response" field.
func ResponseContainsFold(v string) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldContainsFold(FieldResponse, v))
}

// ExecutedAtEQ applies the EQ predicate on the "executed_at" field.
func ExecutedAtEQ(v time.Time) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldEQ(FieldExecutedAt, v))
}

// ExecutedAtNEQ applies the NEQ predicate on the "executed_at" field.
func ExecutedAtNEQ(v time.Time) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldNEQ(FieldExecutedAt, v))
}

// ExecutedAtIn applies the In predicate on the "executed_at" field.
func ExecutedAtIn(vs ...time.Time) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldIn(FieldExecutedAt, vs...))
}

// ExecutedAtNotIn applies the NotIn predicate on the "executed_at" field.
func ExecutedAtNotIn(vs ...time.Time) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldNotIn(FieldExecutedAt, vs...))
}

// ExecutedAtGT applies the GT predicate on the "executed_at" field.
func ExecutedAtGT(v time.Time) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldGT(FieldExecutedAt, v))
}

// ExecutedAtGTE applies the GTE predicate on the "executed_at" field.
func ExecutedAtGTE(v time.Time) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldGTE(FieldExecutedAt, v))
}

// ExecutedAtLT applies the LT predicate on the "executed_at" field.
func ExecutedAtLT(v time.Time) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldLT(FieldExecutedAt, v))
}

// ExecutedAtLTE applies the LTE predicate on the "executed_at" field.
func ExecutedAtLTE(v time.Time) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldLTE(FieldExecutedAt, v))
}

// ExecutedAtIsNil applies the IsNil predicate on the "executed_at" field.
func ExecutedAtIsNil() predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldIsNull(FieldExecutedAt))
}

// ExecutedAtNotNil applies the NotNil predicate on the "executed_at" field.
func ExecutedAtNotNil() predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.FieldNotNull(FieldExecutedAt))
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.CaseyTrigger) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.CaseyTrigger) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.CaseyTrigger) predicate.CaseyTrigger {
	return predicate.CaseyTrigger(sql.NotPredicates(p))
}
