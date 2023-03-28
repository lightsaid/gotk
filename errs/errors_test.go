package errs_test

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/lightsaid/gotk/errs"
	"github.com/stretchr/testify/require"
)

var ErrCustomer = errors.New("customer error")

func TestAsMessage(t *testing.T) {
	err := errs.Success.AsMessage("ok11")
	require.Equal(t, "ok11", err.Message())
}

func openNotExistsFile() error {
	_, err := os.Open("xxx")
	return fmt.Errorf("openNotExistsFile(): %w", err)
}

func TestUnwrap(t *testing.T) {
	fileErr := openNotExistsFile()

	var pe *os.PathError
	require.ErrorAs(t, fileErr, &pe)

	// 第一套娃
	custErr := errs.ServerError.AsException(ErrCustomer)
	require.ErrorIs(t, custErr, ErrCustomer)

	// 第二套娃
	newErr := custErr.AsException(fileErr)
	require.ErrorIs(t, newErr, pe)
	require.ErrorAs(t, newErr, &pe)
}

func TestCustomerAppError(t *testing.T) {
	// 如：业务错误码
	errUserExist := errs.NewAppError(2000, "用户已存在", http.StatusBadRequest)
	require.Equal(t, errUserExist.Code(), 2000)
	require.Equal(t, http.StatusBadRequest, errUserExist.StatusCode())
	require.Equal(t, "用户已存在", errUserExist.Message())

	errUserDel := errs.NewAppError(3000, "用户已删除", http.StatusNotFound)
	require.Equal(t, errUserDel.Code(), 3000)
	require.Equal(t, http.StatusNotFound, errUserDel.StatusCode())
	require.Equal(t, "用户已删除", errUserDel.Message())
}
