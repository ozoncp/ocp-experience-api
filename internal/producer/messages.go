package producer

import (
	"context"

	"github.com/Shopify/sarama"
	"github.com/opentracing/opentracing-go"

	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"

	desc "github.com/ozoncp/ocp-experience-api/pkg/ocp-experience-api"
)

type EventType int

const (
	CreateEvent EventType = iota
	ReadEvent
	UpdateEvent
	DeleteEvent
)

type EventMsg interface {
	sarama.Encoder
}

func NewEvent(ctx context.Context, requestId uint64, eventType EventType, err error) EventMsg {
	e := &event{
		requestId: requestId,
		eventType: eventType,
		err:       err,
	}

	var spanDump opentracing.TextMapCarrier
	span := opentracing.SpanFromContext(ctx)

	if span != nil {
		if err := opentracing.GlobalTracer().Inject(span.Context(), opentracing.TextMap, spanDump); err != nil {
			log.Warn().Msgf("failed to update event message with span info", err)
		} else {
			e.span = spanDump
		}
	}

	return e
}

type event struct {
	requestId   uint64
	eventType   EventType
	err         error
	traceId     string
	span        map[string]string
	encodedData []byte //caching to avoid double encoding on Length() and Encode()
	encodeErr   error
}

func (e *event) Encode() ([]byte, error) {
	if e.encodedData != nil || e.encodeErr != nil {
		return e.encodedData, e.encodeErr
	}

	message := &desc.ExperienceAPIEvent{
		Id: e.requestId,
	}

	if e.err != nil {
		message.Error = e.err.Error()
	}

	switch e.eventType {
	case CreateEvent:
		message.Event = desc.ExperienceAPIEvent_CREATE
	case ReadEvent:
		message.Event = desc.ExperienceAPIEvent_READ
	case UpdateEvent:
		message.Event = desc.ExperienceAPIEvent_UPDATE
	case DeleteEvent:
		message.Event = desc.ExperienceAPIEvent_DELETE
	default:
		log.Panic().Msgf("unexpected event type: %v", e.eventType)
	}

	if len(e.span) > 0 {
		message.TraceSpan = e.span
	}

	e.encodedData, e.encodeErr = proto.Marshal(message)
	return e.encodedData, e.encodeErr
}

func (e *event) Length() int {
	data, _ := e.Encode()
	return len(data)
}
