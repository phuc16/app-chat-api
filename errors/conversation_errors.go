package errors

import "app/pkg/apperror"

const (
	ReadMessageError = 50000 + iota
	UnmarshalMessageError
	AddNewChatError
	GetListUserInConversationError
	CodeConversationNotFound
	NewConversationError
)

func CanNotReadMessage() *apperror.Error {
	return apperror.NewError(ReadMessageError, "can not read message")
}

func CanNotReadUnmarshalMessage() *apperror.Error {
	return apperror.NewError(UnmarshalMessageError, "can not unmarshal message")
}

func CanNotAddNewChat() *apperror.Error {
	return apperror.NewError(AddNewChatError, "can not add new chat")
}

func CanNotGetListIDUserInConversation() *apperror.Error {
	return apperror.NewError(GetListUserInConversationError, "can not get list id user in conversation")
}

func ConversationNotFound() *apperror.Error {
	return apperror.NewError(CodeConversationNotFound, "conversation not found")
}

func DuplicateConversationId() *apperror.Error {
	return apperror.NewError(NewConversationError, "duplicate conversation id")
}
