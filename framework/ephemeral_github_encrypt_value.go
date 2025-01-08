package framework

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/crypto/nacl/box"
)

var _ ephemeral.EphemeralResource = &ephemeralEncryptValue{}

type ephemeralEncryptValue struct{}

func NewEncryptValueEphemeralResource() ephemeral.EphemeralResource {
	return &ephemeralEncryptValue{}
}
func (e *ephemeralEncryptValue) Metadata(ctx context.Context, request ephemeral.MetadataRequest, response *ephemeral.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_encrypt_value" // "github_encrypt_value"
}

func (e *ephemeralEncryptValue) Schema(ctx context.Context, request ephemeral.SchemaRequest, response *ephemeral.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"public_encrypt_key": schema.StringAttribute{
				Required:    true,
				Description: "Public key used to encrypt the value",
			},
			"plaintext_value": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "Plaintext value of the secret to be encrypted.",
			},
			"encrypted_value": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "Plaintext value of the secret to be encrypted.",
			},
		},
	}
}

func (e *ephemeralEncryptValue) Open(ctx context.Context, request ephemeral.OpenRequest, response *ephemeral.OpenResponse) {
	var model encryptValueModel

	diags := request.Config.Get(ctx, &model)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	encryptedBytes, err := encryptPlaintext(model.PlainTextValue.ValueString(), model.PublicEncryptKey.ValueString())
	if err != nil {
		response.Diagnostics.AddError("Unable to encrypt secret",
			fmt.Sprintf("encoutered an error when encrypting the plain text to secret: original error %v", err))
	}
	encryptedValue := base64.StdEncoding.EncodeToString(encryptedBytes)
	model.EncryptedValue = types.StringValue(encryptedValue)

	response.Result.Set(ctx, model)
}

func encryptPlaintext(plaintext, publicKeyB64 string) ([]byte, error) {
	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKeyB64)
	if err != nil {
		return nil, err
	}

	var publicKeyBytes32 [32]byte
	copiedLen := copy(publicKeyBytes32[:], publicKeyBytes)
	if copiedLen == 0 {
		return nil, fmt.Errorf("could not convert publicKey to bytes")
	}

	plaintextBytes := []byte(plaintext)
	var encryptedBytes []byte

	cipherText, err := box.SealAnonymous(encryptedBytes, plaintextBytes, &publicKeyBytes32, nil)
	if err != nil {
		return nil, err
	}

	return cipherText, nil
}

type encryptValueModel struct {
	PublicEncryptKey types.String `tfsdk:"public_encrypt_key"`
	PlainTextValue   types.String `tfsdk:"plaintext_value"`
	EncryptedValue   types.String `tfsdk:"encrypted_value"`
}
