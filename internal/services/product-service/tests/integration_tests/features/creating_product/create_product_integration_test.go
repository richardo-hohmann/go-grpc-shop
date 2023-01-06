package v1

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/rabbitmq"
	consumers2 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/consumers"
	v1 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/creating_product/commands/v1"
	v1_dtos "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/creating_product/dtos/v1"
	v1_event "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/creating_product/events/v1"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/shared/contracts"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/shared/test_fixture/integration"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
	"testing"
)

type createProductIntegrationTests struct {
	*integration.IntegrationTestFixture
}

var consumer *rabbitmq.Consumer

func TestCreateProductIntegration(t *testing.T) {
	suite.Run(t, &createProductIntegrationTests{IntegrationTestFixture: integration.NewIntegrationTestFixture(t, fx.Options(
		fx.Invoke(func(infra *contracts.InfrastructureConfiguration) {
			consumer = rabbitmq.NewConsumer(infra.Cfg.Rabbitmq, infra.ConnRabbitmq, infra.Log, infra.JaegerTracer, consumers2.HandleConsumeCreateProduct)
			err := consumer.ConsumeMessage(infra.Context, v1_event.ProductCreated{})
			if err != nil {
				require.FailNow(t, err.Error())
			}
		}),
	))})
}

func (c *createProductIntegrationTests) Test_Should_Create_New_Product_To_DB() {

	command := v1.NewCreateProduct(gofakeit.Name(), gofakeit.AdjectiveDescriptive(), gofakeit.Price(150, 6000))
	result, err := mediatr.Send[*v1.CreateProduct, *v1_dtos.CreateProductResponseDto](c.Context, command)
	c.Require().NoError(err)

	isPublished := c.RabbitmqPublisher.IsPublished(v1_event.ProductCreated{})
	c.Assert().Equal(true, isPublished)

	isConsumed := consumer.IsConsumed(v1_event.ProductCreated{})
	c.Assert().Equal(true, isConsumed)

	c.Require().NoError(err)

	c.Assert().NotNil(result)
	c.Assert().Equal(command.ProductID, result.ProductId)

	createdProduct, err := c.IntegrationTestFixture.ProductRepository.GetProductById(c.Context, result.ProductId)
	c.Require().NoError(err)
	c.Assert().NotNil(createdProduct)
}

func (c *createProductIntegrationTests) BeforeTest(suiteName, testName string) {
	// some functionality before run tests
}

func (c *createProductIntegrationTests) SetupTest() {
	c.T().Log("SetupTest")
}

func (c *createProductIntegrationTests) TearDownTest() {
	c.T().Log("TearDownTest")
	// cleanup test containers with their hooks
}
