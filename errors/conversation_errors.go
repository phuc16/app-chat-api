package errors

import "app/pkg/apperror"

const (
	ReadMessageError = 50000 + iota
	UnmarshalMessageError
	AddNewChatError
	GetListUserInConversationError
	GetConversationByIdError
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

func CanNotGetConversationById() *apperror.Error {
	return apperror.NewError(GetConversationByIdError, "can not get conversation by id")
}

func DuplicateConversationId() *apperror.Error {
	return apperror.NewError(NewConversationError, "duplicate conversation id")
}
