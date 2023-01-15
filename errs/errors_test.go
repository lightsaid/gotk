package errs_test

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/gotk/errs"
	"github.com/stretchr/testify/require"
)

func TestAsMessage(t *testing.T) {
	err := errs.InvalidParams.AsMessage("customer msg")
	require.Equal(t, "customer msg", err.Message())
	require.Equal(t, "参数错误", errs.InvalidParams.Message())
}

func openNotExistsFile() error {
	_, err := os.Open("xxx")
	return fmt.Errorf("openNotExistsFile(): %w", err)
}

var ErrCustomer = errors.New("customer error")

func TestAsException(t *testing.T) {
	fileErr := openNotExistsFile()

	var pe *os.PathError
	log.Println("err-> As", errors.As(fileErr, &pe)) // true
	require.ErrorAs(t, fileErr, &pe)

	// 第一套娃
	custErr := errs.ServerError.AsException(ErrCustomer)

	log.Println("custErr-> Is ", errors.Is(custErr, ErrCustomer)) // true
	log.Println(custErr)

	require.ErrorIs(t, custErr, ErrCustomer)

	// 第二套娃
	newErr := custErr.AsException(fileErr)

	log.Println(newErr)

	log.Println("newErr-> Is ", errors.Is(newErr, pe))
	require.ErrorIs(t, newErr, pe)

	log.Println("newErr-> As ", errors.As(newErr, &pe))
	require.ErrorAs(t, newErr, &pe)
}

// 测试自定义错误
func TestCustomerApperror(t *testing.T) {
	// 如：业务错误码
	errUserExist := errs.NewAppError(2000, "用户已存在")

	require.Equal(t, errUserExist.Code(), 2000)

	require.Equal(t, http.StatusInternalServerError, errUserExist.StatusCode())

	errUserDel := errs.NewAppError(2001, "用户已删除", http.StatusNotFound)

	require.Equal(t, http.StatusNotFound, errUserDel.StatusCode())
}
