package detector

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type DetectorClient struct {
	BaseClient
}

func NewDetectorClient(subscriptionID string) DetectorClient {
	return NewDetectorClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

func NewDetectorClientWithBaseURI(baseURI string, subscriptionID string) DetectorClient {
	return DetectorClient{NewWithBaseURI(baseURI, subscriptionID)}
}

func (client DetectorClient) Get(ctx context.Context, resourceGroupName string, resourceName string, detectId string, startTime string, endTime string) (result Detector, err error) {

	req, err := client.GetPreparer(ctx, resourceGroupName, resourceName, detectId, startTime, endTime)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerservice.detector", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "containerservice.detector", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerservice.detector", "Get", resp, "Failure responding to request")
		return
	}

	return

}

func (client DetectorClient) GetPreparer(ctx context.Context, resourceGroupName string, resourceName string, detectId string, startTime string, endTime string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"resourceName":      autorest.Encode("path", resourceName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"detectId":          autorest.Encode("path", detectId),
	}

	const APIVersion = "2019-04-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
		"startTime":   startTime,
		"endTime":     endTime,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerService/managedClusters/{resourceName}/detectors/{detectId}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

func (client DetectorClient) GetSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

func (client DetectorClient) GetResponder(resp *http.Response) (result Detector, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
