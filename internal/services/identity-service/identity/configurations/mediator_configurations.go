package configurations

import (
	"github.com/mehdihadeli/go-mediatr"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/identity/dtos"
	v1 "github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/identity/features/registering_user/commands/v1"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/shared/contracts"
)

func ConfigUsersMediator(infra *contracts.InfrastructureConfiguration) error {

	//https://stackoverflow.com/questions/72034479/how-to-implement-generic-interfaces
	err := mediatr.RegisterRequestHandler[*v1.RegisterUser, *dtos.RegisterUserResponseDto](v1.NewRegisterUserHandler(infra))
	if err != nil {
		return err
	}

	return nil
}
