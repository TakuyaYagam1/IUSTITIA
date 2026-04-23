// Package domain holds IUSTITIA business entities. One file per aggregate.
// Domain types are free of persistence and transport concerns - repo
// converts rows here, controller converts these to openapi payloads.
package domain

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleCitizen    Role = "citizen"
	RoleProsecutor Role = "prosecutor"
	RoleJudge      Role = "judge"
	RoleRegistrar  Role = "registrar"
)

func (r Role) String() string { return string(r) }

func (r Role) Valid() bool {
	switch r {
	case RoleCitizen, RoleProsecutor, RoleJudge, RoleRegistrar:
		return true
	}
	return false
}

type User struct {
	ID        uuid.UUID
	Username  string
	Password  string
	Role      Role
	Dome      string
	CreatedAt time.Time
}
