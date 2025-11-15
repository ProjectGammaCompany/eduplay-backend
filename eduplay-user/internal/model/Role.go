package model

import (
	dto "eduplay-user/internal/generated"
)

type Role string

const (
	USER      Role = Role(dto.Role_USER)
	OPERATOR  Role = Role(dto.Role_OPERATOR)
	MODERATOR Role = Role(dto.Role_MODERATOR)
	JURIST    Role = Role(dto.Role_JURIST)
)

// var roleMap = map[dto.Role]Role{
// 	dto.Role_USER:      USER,
// 	dto.Role_OPERATOR:  OPERATOR,
// 	dto.Role_MODERATOR: MODERATOR,
// 	dto.Role_JURIST:    JURIST,
// }

// func ToModelRole(protoRole dto.Role) Role {
// 	return roleMap[protoRole]
// }
