package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/http/echo/middleware"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/utils"
	v1 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/getting_products/dtos/v1"
	query_v1 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/getting_products/queries/v1"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/shared"
	"net/http"
)

type getProductsEndpoint struct {
	*shared.ProductEndpointBase[shared.InfrastructureConfiguration]
}

func NewGetProductsEndpoint(productEndpointBase *shared.ProductEndpointBase[shared.InfrastructureConfiguration]) *getProductsEndpoint {
	return &getProductsEndpoint{productEndpointBase}
}

func (ep *getProductsEndpoint) MapRoute() {
	ep.ProductsGroup.GET("", ep.getAllProducts(), middleware.ValidateBearerToken())
}

// GetAllProducts
// @Tags Products
// @Summary Get all product
// @Description Get all products
// @Accept json
// @Produce json
// @Param GetProductsRequestDto query v1.GetProductsRequestDto false "GetProductsRequestDto"
// @Success 200 {object} v1.GetProductsResponseDto
// @Security ApiKeyAuth
// @Router /api/v1/products [get]
func (ep *getProductsEndpoint) getAllProducts() echo.HandlerFunc {
	return func(c echo.Context) error {

		ctx := c.Request().Context()

		listQuery, err := utils.GetListQueryFromCtx(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		request := &v1.GetProductsRequestDto{ListQuery: listQuery}
		if err := c.Bind(request); err != nil {
			ep.Configuration.Log.Warn("Bind", err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		query := query_v1.NewGetProducts(request.ListQuery)

		queryResult, err := mediatr.Send[*query_v1.GetProducts, *v1.GetProductsResponseDto](ctx, query)

		if err != nil {
			ep.Configuration.Log.Warnf("GetProducts", err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		return c.JSON(http.StatusOK, queryResult)
	}
}
