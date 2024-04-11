package infra

import (
	"net"
	"time"

	"github.com/segmentio/kafka-go"
)

// KafkaWriter is a wrapper of kafka writer.
type KafkaWriter struct {
	*kafka.Writer
	Coco
}

// KafkaReader is a wrapper of kafka reader.
type KafkaReader struct {
	*kafka.Reader
	Coco
}

type KafkaMessage = kafka.Message

type KafkaWriterConfig struct {
	Brokers      []string
	Topic        string
	ClientID     string
	DialTimeout  time.Duration
	BatchTimeout time.Duration
	Async        bool
}

type KafkaReaderConfig struct {
	Brokers     []string
	Topic       string
	GroupID     string
	GroupTopics []string
}

// NewKafkaWriter creates a new kafka writer.
func NewKafkaWriter(c KafkaWriterConfig) (*KafkaWriter, error) {
	w := &kafka.Writer{
		Transport: &kafka.Transport{
			ClientID: c.ClientID,
			Dial: (&net.Dialer{
				Timeout: c.DialTimeout,
			}).DialContext,
		},
		Addr:         kafka.TCP(c.Brokers...),
		Async:        c.Async,
		BatchTimeout: c.BatchTimeout,
		Topic:        c.Topic,
		RequiredAcks: kafka.RequireNone,
	}
	return &KafkaWriter{Writer: w}, nil
}

// NewKafkaReader creates a new kafka reader.
func NewKafkaReader(c KafkaReaderConfig) (*KafkaReader, error) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     c.Brokers,
		Topic:       c.Topic,
		GroupID:     c.GroupID,
		GroupTopics: c.GroupTopics,
	})
	return &KafkaReader{Reader: r}, nil
}
