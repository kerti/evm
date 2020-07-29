package functional

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/kerti/evm/02-kitara-store/model"
	"github.com/stretchr/testify/assert"
)

const (
	url = "http://localhost:8080/orders/process"
)

func TestOrderConcurrency(t *testing.T) {

	payload1, _ := json.Marshal(model.OrderProcessInput{
		OrderID: uuid.FromStringOrNil("5b27773a-efce-4e21-8474-2694ebdaa084"),
	})
	body1 := strings.NewReader(string(payload1))

	payload2, _ := json.Marshal(model.OrderProcessInput{
		OrderID: uuid.FromStringOrNil("63b2ae3c-040a-4d64-89b2-b5e7e1951e89"),
	})
	body2 := strings.NewReader(string(payload2))

	resp1 := make(chan *http.Response)
	resp2 := make(chan *http.Response)

	// Run both requests concurrently
	go func(t *testing.T) {
		r1, err := http.Post(url, "application/json", body1)
		assert.Nil(t, err)
		resp1 <- r1
	}(t)

	go func(t *testing.T) {
		r2, err := http.Post(url, "application/json", body2)
		assert.Nil(t, err)
		resp2 <- r2
	}(t)

	response1 := <-resp1
	response2 := <-resp2

	// If one of the responses returns OK, the other one should not return OK
	if response1.StatusCode == http.StatusOK {
		assert.NotEqual(t, response2.StatusCode, http.StatusOK)
	}

	if response2.StatusCode == http.StatusOK {
		assert.NotEqual(t, response1.StatusCode, http.StatusOK)
	}

}
