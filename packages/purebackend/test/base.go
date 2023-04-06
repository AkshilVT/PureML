package test

import (
	uuid "github.com/satori/go.uuid"
)

const (
	ValidAdminToken  = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImRlbW9AYXp0bGFuLmluIiwiaGFuZGxlIjoiZGVtbyIsInV1aWQiOiIxMTExMTExMS0xMTExLTExMTEtMTExMS0xMTExMTExMTExMTEifQ.dpM9Ij_Y25A5yNiVTt8hI-ZtjDqUfvbAFdtU9-RyDbs"
	ValidUserToken   = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Im5vdGFkbWluQGF6dGxhbi5pbiIsImhhbmRsZSI6Im5vdGFkbWluIiwidXVpZCI6IjIyMjIyMjIyLTIyMjItMjIyMi0yMjIyLTIyMjIyMjIyMjIyMiJ9.H_m-iWN4M1aZNZ3CB914kFm3JkQzJRy5x9uc6a0KE9c"
	ValidTokenNoUser = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Im5vb25lQG5vb25lLm9uZSIsImhhbmRsZSI6Im5vcGUiLCJ1dWlkIjoiMTExMTExMTEtMjIyMi0zMzMzLTQ0NDQtMTExMTExMTExMTExIn0.yxhQWKYh6TW6DpfaYHM0LvDNKGSoy080jGvyMg23wYQ"
	InvalidToken     = "Bearer abcdefghijklmnopqrstuvwxyz"
)

// Valid		valid uuid format
// AdminUser	admin user
// User			normal user
// NoUser		user not found
// NoOrg		org not found
var (
	ValidAdminUserUuid    = uuid.Must(uuid.FromString("11111111-1111-1111-1111-111111111111"))
	ValidUserUuid         = uuid.Must(uuid.FromString("22222222-2222-2222-2222-222222222222"))
	ValidNoUserUuid       = uuid.Must(uuid.FromString("11111111-2222-3333-4444-111111111111"))
	InvalidOrgUuidString  = "11111111-2222-3333-4444-11111111"
	ValidNoOrgUuid        = uuid.Must(uuid.FromString("11111111-2222-3333-4444-111111111111"))
	ValidAdminUserOrgUuid = uuid.Must(uuid.FromString("11111111-1111-1111-1111-111111111111"))
	ValidUserOrgUuid      = uuid.Must(uuid.FromString("22222222-2222-2222-2222-222222222222"))
	ValidSessionUuid      = uuid.Must(uuid.FromString("11111111-1111-1111-1111-111111111111"))
	ValidNoSessionUuid    = uuid.Must(uuid.FromString("11111111-2222-3333-4444-111111111111"))
)
