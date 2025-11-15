package converters

import (
	dto "eduplay-user/internal/generated"
	// "user-service/internal/model"
)

// func RoleToDto(role model.Role) dto.Role {
// 	switch role {
// 	case model.USER:
// 		return dto.Role_USER
// 	case model.OPERATOR:
// 		return dto.Role_OPERATOR
// 	case model.MODERATOR:
// 		return dto.Role_MODERATOR
// 	case model.JURIST:
// 		return dto.Role_JURIST
// 	default:
// 		// Возвращаем значение по умолчанию или обрабатываем ошибку
// 		return dto.Role_USER
// 	}
// }

func StringToDto(role string) dto.Role {
	switch role {
	case "user":
		return dto.Role_USER
	case "operator":
		return dto.Role_OPERATOR
	case "moderator":
		return dto.Role_MODERATOR
	case "jurist":
		return dto.Role_JURIST
	default:
		// Возвращаем значение по умолчанию или обрабатываем ошибку
		return dto.Role_USER
	}
}
