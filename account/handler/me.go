package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"memrizr/account/observability/logger"
	"net/http"
)

// Me handler calls services for getting
// a user's details
func (h *Handler) Me(c *gin.Context) {
	//// A *model.User will eventually be added to context in middleware
	//user, exists := c.Get("user")
	//
	//// This shouldn't happen, as our middleware ought to throw an error.
	//// This is an extra safety measure
	//// We'll extract this logic later as it will be common to all handler
	//// methods which require a valid user
	//if !exists {
	//	log.Printf("Unable to extract user from request context for unknown reason: %v\n", c)
	//	err := apperrors.NewInternal()
	//	c.JSON(err.Status(), gin.H{
	//		"error": err,
	//	})
	//
	//	return
	//}
	//
	//uid := user.(*model.User).UID
	//
	//// gin.Context satisfies go's context.Context interface
	//u, err := h.UserService.Get(c, uid)
	//
	//if err != nil {
	//	log.Printf("Unable to find user: %v\n%v", uid, err)
	//	e := apperrors.NewNotFound("user", uid.String())
	//
	//	c.JSON(e.Status(), gin.H{
	//		"error": e,
	//	})
	//	return
	//}
	//
	//c.JSON(http.StatusOK, gin.H{
	//	"user": u,
	//})
	var tracer = otel.GetTracerProvider().Tracer("me")
	_, span := tracer.Start(c.Request.Context(), "ME")
	defer span.End()

	logger.LogInfo("Get Me successfully", logrus.Fields{
		"Hi": "success",
	})
	c.JSON(http.StatusOK, gin.H{
		"message": "hello world is it save or not",
	})
}
