package ent

//go:generate go run -mod=mod entgo.io/ent/cmd/ent generate --template=./schema/templates --feature intercept,schema/snapshot,sql/upsert ./schema
