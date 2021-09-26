package lib

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
)

func ParseGinRequestToProtobuf(ctx *gin.Context , m proto.Message) error{
	payload, err := ioutil.ReadAll(ctx.Request.Body)
	// parse with protojson
	if ( err != nil){
		return err
	}
	err = protojson.Unmarshal(payload,m)
	if ( err != nil){
		return err
	}
	return nil
}
