package inalog

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

type testCtxKey string

func TestFiberContext(t *testing.T) {
	assert := assert.New(t)
	app := fiber.New()
	app.Get("/", func(f *fiber.Ctx) error {
		ctx0 := context.WithValue(f.UserContext(), testCtxKey("trace-id"), "1234")
		f.SetUserContext(ctx0)
		f.Context().SetUserValue(testCtxKey("span-id"), "4445")
		ctx := context.WithValue(WithFiberCtx(f.Context()), testCtxKey("subspan-id"), "0090")

		assert.Equal("4445", ctx.Value(testCtxKey("span-id")))
		assert.Equal("1234", ctx.Value(testCtxKey("trace-id")))
		assert.Equal("0090", ctx.Value(testCtxKey("subspan-id")))
		fmt.Println(ctx)
		return nil
	})

	req := httptest.NewRequest(fiber.MethodGet, "/", nil)

	_, err := app.Test(req)
	assert.NoError(err)
}
