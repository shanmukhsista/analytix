package controller

import "analytix/generated"

type EventAcquisitionAPI interface {
	 AcquireEvents(eventPayload * generated.AcquisitionEvent) (*generated.AcquisitionEventResponse,error)
}
