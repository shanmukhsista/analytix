package service

import (
	"analytix/generated"
	"analytix/pkg/controller"
	"analytix/pkg/lib"
	"github.com/ansel1/merry/v2"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AcquisitionRestGinService struct {
	// Connect to a controller who can handle API Requests.
	acquisitionController controller.EventAcquisitionAPI
}

func (a AcquisitionRestGinService) GetEvents(ctx *gin.Context) {
	// Check if we can
	// use headers if possible to restrict scope.
	// Check if there is a valid client id .
	// parse with protojson
	var eventPayload generated.AcquisitionEvent
	err := lib.ParseGinRequestToProtobuf(ctx, &eventPayload)
	if ( err != nil){
		err := merry.New("An error occoured while parsing the request. ",merry.WithHTTPCode(http.StatusBadRequest))
		ctx.Error(err)
		return
	}
	// Handle error.
	response, _ := a.acquisitionController.AcquireEvents(&eventPayload)
	// check if we have 1 or more than 1 event.
	ctx.JSON(http.StatusAccepted,response)
}

func NewRestAcquisitionService(ac controller.EventAcquisitionAPI) *AcquisitionRestGinService {
	return &AcquisitionRestGinService{acquisitionController: ac}
}