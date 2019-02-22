package facade

import (
	"tianwei.pro/sam/core/dto"
)

type TokenFacade interface {
	DecodeToken(token string)(*dto.UserDto, error)
	EncodeToken(user *dto.UserDto, Terminal int8) (string, error)
}