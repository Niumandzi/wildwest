package errors

import (
	"fmt"
	"wildwest/pkg/contextutils"
)

type RepoError struct {
	UserID        string
	OperationName string
	OperationID   string
	Message       string
}

func (e *RepoError) Error() string {
	return fmt.Sprintf("Operation Name %s, Operation ID %s, User ID %s : %s", e.OperationName, e.OperationID, e.UserID, e.Message)
}

func NewRepoError(contextData contextutils.ContextData, message string) error {
	return &RepoError{
		OperationName: contextData.OperationName,
		OperationID:   contextData.OperationID,
		UserID:        contextData.UserID,
		Message:       message,
	}
}

// CreateError Ошибка создания записи
func CreateError(contextData contextutils.ContextData, entity string, err error) error {
	return NewRepoError(contextData, fmt.Sprintf("failed to create %s: %v", entity, err))
}

// RecordNotFoundError Ошибка нахождения записи для обновления
func RecordNotFoundError(contextData contextutils.ContextData, entity string) error {
	return NewRepoError(contextData, fmt.Sprintf("no record found %s", entity))
}

// UpdateError Ошибка обновления записи
func UpdateError(contextData contextutils.ContextData, entity string, err error) error {
	return NewRepoError(contextData, fmt.Sprintf("failed to update %s: %v", entity, err))
}

// DeleteError Ошибка удаления записи
func DeleteError(contextData contextutils.ContextData, entity string, err error) error {
	return NewRepoError(contextData, fmt.Sprintf("failed to delete %s: %v", entity, err))
}

// TransactionStartError Ошибка начала транзакции
func TransactionStartError(contextData contextutils.ContextData, err error) error {
	return NewRepoError(contextData, fmt.Sprintf("error starting transaction: %v", err))
}

// TransactionCommitError Ошибка завершения транзакции
func TransactionCommitError(contextData contextutils.ContextData, err error) error {
	return NewRepoError(contextData, fmt.Sprintf("transaction commit failed: %v", err))
}
