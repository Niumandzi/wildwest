package contextutils

import (
	"context"
	"github.com/google/uuid"
	"net/http"
)

type ContextData struct {
	OperationName string
	OperationID   string
	UserID        string
}

func NewContext(r *http.Request, userID int, operationName string) context.Context {
	operationID := uuid.New().String()
	ctx := context.WithValue(r.Context(), "operationName", operationName)
	ctx = context.WithValue(ctx, "operationID", operationID)
	ctx = context.WithValue(ctx, "userID", userID)
	return ctx
}

func ExtractContextData(ctx context.Context) ContextData {
	userID, _ := ctx.Value("userID").(string)
	opName, _ := ctx.Value("operationName").(string)
	opID, _ := ctx.Value("operationID").(string)

	if opName == "" {
		opName = "unknown"
	}
	if opID == "" {
		opID = "unknown"
	}
	if userID == "" {
		userID = "unknown"
	}

	return ContextData{
		OperationName: opName,
		OperationID:   opID,
		UserID:        userID,
	}
}
