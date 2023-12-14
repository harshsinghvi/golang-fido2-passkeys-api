package autoroutes

import (
	"github.com/harshsinghvi/golang-fido2-passkeys-api/autoroutes/controllers"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/autoroutes/models"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/autoroutes/routes"
)

// Exports

type Routes = models.Routes
type Route = models.Route
type Config = models.Config
type GenerateFunction = models.GenerateFunction
type GenerateFields = models.GenerateFields
type ValidationFunction = models.ValidationFunction
type ValidationFields = models.ValidationFields

const (
	MethodGet    = models.MethodGet
	MethodPost   = models.MethodPost
	MethodPut    = models.MethodPut
	MethodDelete = models.MethodDelete
)

var (
	GetController    = controllers.GetController
	PostController   = controllers.PostController
	DeleteController = controllers.DeleteController
	PutController    = controllers.PutController
)

var GenerateRoutes = routes.GenerateRoutes
