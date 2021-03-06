package worker

import (
	"encoding/json"
	"errors"
	"github.com/dreamvo/gilfoyle/worker/mocks"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProducers(t *testing.T) {
	t.Run("ProduceVideoTranscodingQueue", func(t *testing.T) {
		t.Run("should publish a new message", func(t *testing.T) {
			params := VideoTranscodingParams{
				MediaUUID:      uuid.New(),
				SourceFilePath: "test",
			}

			body, err := json.Marshal(params)
			assert.NoError(t, err)

			ch := new(mocks.MockedChannel)

			ch.On("Publish", "", VideoTranscodingQueue, false, false, amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "application/json",
				Body:         body,
			}).Return(nil)

			err = ProduceVideoTranscodingQueue(ch, params)
			assert.NoError(t, err)

			ch.AssertExpectations(t)
		})

		t.Run("should publish a new message with AMQP error", func(t *testing.T) {
			params := VideoTranscodingParams{
				MediaUUID:      uuid.New(),
				SourceFilePath: "test",
			}

			body, err := json.Marshal(params)
			assert.NoError(t, err)

			ch := new(mocks.MockedChannel)

			ch.On("Publish", "", VideoTranscodingQueue, false, false, amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "application/json",
				Body:         body,
			}).Return(errors.New("test"))

			err = ProduceVideoTranscodingQueue(ch, params)
			assert.EqualError(t, err, "test")

			ch.AssertExpectations(t)
		})
	})
}
