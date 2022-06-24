package provider

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/crypto/nacl/box"
)

func dataSourceSodiumEncryptedItem() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSodiumEncryptedItemRead,

		Description: "Encrypt a string value with given public key using sodium library",

		Schema: map[string]*schema.Schema{
			"public_key": {
				Type:        schema.TypeString,
				Description: "Public key to use when encrypting",
				Required:    true,
				ForceNew:    true,
			},
			"content_base64": {
				Type:        schema.TypeString,
				Description: "Base64 encoded version of the raw string to encrypt.",
				Sensitive:   true,
				Required:    true,
				ForceNew:    true,
				StateFunc:   dataSourceSodiumEncryptedItemState,
			},
			"encrypted_value_base64": {
				Type:        schema.TypeString,
				Description: "Base64 encoded version of the encrypted result .",
				Sensitive:   true,
				Computed:    true,
			},
		},
	}
}

func dataSourceSodiumEncryptedItemRead(d *schema.ResourceData, m interface{}) error {

	secret, err := base64.StdEncoding.DecodeString(d.Get("content_base64").(string))
	if err != nil {
		return fmt.Errorf("failed to decode base64 content")
	}
	secretBytes := []byte(secret)

	// Getting public key from input
	var pkBytes [32]byte
	copy(pkBytes[:], []byte(d.Get("public_key").(string)))

	// Encrypting string with given pubKey
	enc, err := box.SealAnonymous(nil, secretBytes, &pkBytes, nil)

	// Encoding result to base64
	encEnc := base64.StdEncoding.EncodeToString(enc)
	d.Set("encrypted_value_base64", encEnc)

	checksum := sha1.Sum([]byte(encEnc))
	d.SetId(hex.EncodeToString(checksum[:]))

	return nil
}

// Store a hash of an input instead of an actual one
func dataSourceSodiumEncryptedItemState(i interface{}) string {
	checksum := sha1.Sum([]byte((i).(string)))
	return hex.EncodeToString(checksum[:])
}
