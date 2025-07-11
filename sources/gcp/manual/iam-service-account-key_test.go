package manual_test

import (
	"context"
	"fmt"
	"testing"

	"cloud.google.com/go/iam/admin/apiv1/adminpb"
	"go.uber.org/mock/gomock"

	"github.com/overmindtech/cli/sdp-go"
	"github.com/overmindtech/cli/sources"
	"github.com/overmindtech/cli/sources/gcp/manual"
	gcpshared "github.com/overmindtech/cli/sources/gcp/shared"
	"github.com/overmindtech/cli/sources/gcp/shared/mocks"
	"github.com/overmindtech/cli/sources/shared"
)

func TestIAMServiceAccountKey(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockIAMServiceAccountKeyClient(ctrl)
	projectID := "test-project-id"

	testServiceAccount := "test-sa@test-project-id.iam.gserviceaccount.com"
	testKeyName := "1234567890abcdef"
	testKeyFullName := "projects/test-project-id/serviceAccounts/test-sa@test-project-id.iam.gserviceaccount.com/keys/1234567890abcdef"

	t.Run("Get", func(t *testing.T) {
		wrapper := manual.NewIAMServiceAccountKey(mockClient, projectID)

		mockClient.EXPECT().Get(ctx, gomock.Any()).Return(createServiceAccountKey(testKeyFullName), nil)

		adapter := sources.WrapperToAdapter(wrapper)

		sdpItem, qErr := adapter.Get(ctx, wrapper.Scopes()[0], shared.CompositeLookupKey(testServiceAccount, testKeyName), true)
		if qErr != nil {
			t.Fatalf("Expected no error, got: %v", qErr)
		}

		t.Run("StaticTests", func(t *testing.T) {
			queryTests := shared.QueryTests{
				{
					ExpectedType:             gcpshared.IAMServiceAccount.String(),
					ExpectedMethod:           sdp.QueryMethod_GET,
					ExpectedQuery:            testServiceAccount,
					ExpectedScope:            projectID,
					ExpectedBlastPropagation: &sdp.BlastPropagation{In: true, Out: false},
				},
			}

			shared.RunStaticTests(t, adapter, sdpItem, queryTests)
		})
	})

	t.Run("Search", func(t *testing.T) {
		wrapper := manual.NewIAMServiceAccountKey(mockClient, projectID)
		adapter := sources.WrapperToAdapter(wrapper)

		mockClient.EXPECT().Search(ctx, gomock.Any()).Return(&adminpb.ListServiceAccountKeysResponse{
			Keys: []*adminpb.ServiceAccountKey{
				createServiceAccountKey(testKeyFullName),
			},
		}, nil)

		sdpItems, err := adapter.Search(ctx, wrapper.Scopes()[0], testServiceAccount, true)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		expectedCount := 1
		actualCount := len(sdpItems)
		if actualCount != expectedCount {
			t.Fatalf("Expected %d items, got: %d", expectedCount, actualCount)
		}

		for _, item := range sdpItems {
			if err := item.Validate(); err != nil {
				t.Fatalf("Expected no validation error, got: %v", err)
			}
		}
	})

	t.Run("SearchWithTerraformQueryMap", func(t *testing.T) {
		wrapper := manual.NewIAMServiceAccountKey(mockClient, projectID)

		mockClient.EXPECT().Get(ctx, gomock.Any()).Return(createServiceAccountKey(testKeyFullName), nil)

		adapter := sources.WrapperToAdapter(wrapper)

		// projects/{{project}}/serviceAccounts/{{account}}/keys/{{key}}
		terraformResourceID := fmt.Sprintf("projects/%s/serviceAccounts/%s/keys/%s", projectID, testServiceAccount, testKeyName)

		sdpItems, err := adapter.Search(ctx, wrapper.Scopes()[0], terraformResourceID, true)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		expectedCount := 1
		actualCount := len(sdpItems)
		if actualCount != expectedCount {
			t.Fatalf("Expected %d items, got: %d", expectedCount, actualCount)
		}

		if err := sdpItems[0].Validate(); err != nil {
			t.Fatalf("Expected no validation error, got: %v", err)
		}
	})

}

// createServiceAccountKey creates a ServiceAccountKey with the specified name.
func createServiceAccountKey(name string) *adminpb.ServiceAccountKey {
	return &adminpb.ServiceAccountKey{
		Name: name,
	}
}
