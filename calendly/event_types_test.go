package calendly

import (
	"github.com/stretchr/testify/assert"
	"fmt"
	"net/http"
	"context"
)

func (suite *CalendlyClientTestSuite) TestEventTypesService_ListEventTypes() {
	assert := assert.New(suite.T())
	route := fmt.Sprintf("/%s", eventTypesPath)

	suite.mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, http.MethodGet)
		fmt.Fprint(w, `{"data":[{"id":"123"},{"id":"456"}]}`)
	})

	eventTypes, _, err := suite.client.EventTypes.List(context.Background(), nil)
	assert.Nil(err)

	want := []*EventType{
		{ID: "123"},
		{ID: "456"},
	}
	assert.Equal(want, eventTypes)
}
