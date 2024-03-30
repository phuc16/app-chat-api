package service

//go:generate mockgen -source $GOFILE -destination ../mocks/$GOPACKAGE/mock_$GOFILE -package mocks

type IUserRepo interface {
}
