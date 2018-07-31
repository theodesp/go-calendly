package calendly

import (
	"github.com/stretchr/testify/assert"
	"fmt"
	"net/http"
	"context"
)

func (suite *CalendlyClientTestSuite) TestUsersService_AboutMe() {
	assert := assert.New(suite.T())
	route := fmt.Sprintf("/%s", aboutMePath)

	suite.mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, http.MethodGet)
		fmt.Fprint(w, `{"data":{"id":"123"}}`)
	})

	me, _, err := suite.client.Users.AboutMe(context.Background())
	assert.Nil(err)

	want := &AboutMe{ID: "123"}
	assert.Equal(want, me)
}
