package shared

import (
	"context"

	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
	"github.com/googleapis/gax-go/v2"
)

// CloudKMSKeyRingIterator is an interface for iterating over KMS KeyRings
type CloudKMSKeyRingIterator interface {
	Next() (*kmspb.KeyRing, error)
}

// CloudKMSKeyRingClient is an interface for the KMS KeyRing client
type CloudKMSKeyRingClient interface {
	Get(ctx context.Context, req *kmspb.GetKeyRingRequest, opts ...gax.CallOption) (*kmspb.KeyRing, error)
	Search(ctx context.Context, req *kmspb.ListKeyRingsRequest, opts ...gax.CallOption) CloudKMSKeyRingIterator
}

// cloudKMSKeyRingClient is a concrete implementation of CloudKMSKeyRingClient
type cloudKMSKeyRingClient struct {
	client *kms.KeyManagementClient
}

// NewCloudKMSKeyRingClient creates a new CloudKMSKeyRingClient
func NewCloudKMSKeyRingClient(keyRingClient *kms.KeyManagementClient) CloudKMSKeyRingClient {
	return &cloudKMSKeyRingClient{
		client: keyRingClient,
	}
}

// Get retrieves a KMS KeyRing
func (c cloudKMSKeyRingClient) Get(ctx context.Context, req *kmspb.GetKeyRingRequest, opts ...gax.CallOption) (*kmspb.KeyRing, error) {
	return c.client.GetKeyRing(ctx, req, opts...)
}

// List lists KMS KeyRings and returns an iterator
func (c cloudKMSKeyRingClient) Search(ctx context.Context, req *kmspb.ListKeyRingsRequest, opts ...gax.CallOption) CloudKMSKeyRingIterator {
	return c.client.ListKeyRings(ctx, req, opts...)
}

// CloudKMSCryptoKeyVersionIterator is an interface for iterating over Cloud KMS CryptoKeyVersions
type CloudKMSCryptoKeyVersionIterator interface {
	Next() (*kmspb.CryptoKeyVersion, error)
}

// CloudKMSCryptoKeyVersionClient is an interface for the Cloud KMS CryptoKeyVersion client
type CloudKMSCryptoKeyVersionClient interface {
	Get(ctx context.Context, req *kmspb.GetCryptoKeyVersionRequest, opts ...gax.CallOption) (*kmspb.CryptoKeyVersion, error)
	List(ctx context.Context, req *kmspb.ListCryptoKeyVersionsRequest, opts ...gax.CallOption) CloudKMSCryptoKeyVersionIterator
}

// CloudKMSCryptoKeyIterator is an interface for iterating over KMS CryptoKeys
type CloudKMSCryptoKeyIterator interface {
	Next() (*kmspb.CryptoKey, error)
}

// CloudKMSCryptoKeyClient is an interface for the KMS CryptoKey client
type CloudKMSCryptoKeyClient interface {
	Get(ctx context.Context, req *kmspb.GetCryptoKeyRequest, opts ...gax.CallOption) (*kmspb.CryptoKey, error)
	List(ctx context.Context, req *kmspb.ListCryptoKeysRequest, opts ...gax.CallOption) CloudKMSCryptoKeyIterator
}

// cloudKMSCryptoKeyClient is a concrete implementation of CloudKMSCryptoKeyClient
type cloudKMSCryptoKeyClient struct {
	client *kms.KeyManagementClient
}

// NewCloudKMSCryptoKeyClient creates a new CloudKMSCryptoKeyClient
func NewCloudKMSCryptoKeyClient(cryptoKeyClient *kms.KeyManagementClient) CloudKMSCryptoKeyClient {
	return &cloudKMSCryptoKeyClient{
		client: cryptoKeyClient,
	}
}

// Get retrieves a KMS CryptoKey
func (c cloudKMSCryptoKeyClient) Get(ctx context.Context, req *kmspb.GetCryptoKeyRequest, opts ...gax.CallOption) (*kmspb.CryptoKey, error) {
	// Client options are ignored on individual calls
	return c.client.GetCryptoKey(ctx, req, opts...)
}

// List lists KMS CryptoKeys and returns an iterator
func (c cloudKMSCryptoKeyClient) List(ctx context.Context, req *kmspb.ListCryptoKeysRequest, opts ...gax.CallOption) CloudKMSCryptoKeyIterator {
	return c.client.ListCryptoKeys(ctx, req, opts...)
}
