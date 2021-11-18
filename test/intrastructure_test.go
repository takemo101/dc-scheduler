package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/takemo101/dc-scheduler/app"
	"github.com/takemo101/dc-scheduler/boot"
	"github.com/takemo101/dc-scheduler/pkg"
	"github.com/takemo101/dc-scheduler/pkg/domain"
	"go.uber.org/fx"
)

func Test_Infrastructure(t *testing.T) {
	boot.Testing(
		t,
		boot.TestOptions{
			ConfigPath:       "../config.testing.yml",
			CurrentDirectory: "../",
			FXOption: fx.Options(
				app.Module,
				pkg.Module,
			),
		},
		func(
			crypter domain.UserActivationSignatureCrypter,
		) {
			t.Run("UserActivationSignatureCrypter", func(t *testing.T) {
				key, err := domain.GenerateUserActivationKey()
				assert.Equal(t, nil, err)

				userSignature := domain.UserSignature{
					Email:         domain.UserEmail("test@test.com"),
					ActivationKey: key,
				}

				// 暗号化
				signature, err := crypter.Encrypt(userSignature)
				assert.Equal(t, nil, err)

				// 複合化
				decryptSignature, err := crypter.Decrypt(signature)
				assert.Equal(t, nil, err)

				assert.True(
					t,
					userSignature.ActivationKey.Equals(
						decryptSignature.ActivationKey,
					),
				)

				assert.True(
					t,
					userSignature.Email.Equals(
						decryptSignature.Email,
					),
				)
			})
		},
	)
}
