package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Feinot/metric-and-allert/internal/forms"
	"github.com/Feinot/metric-and-allert/internal/logger"
)

var ServerCounter = make(map[string]int64)
var ServerGauge = make(map[string]float64)
var AgentCounter = make(map[string]int64)
var AgentGauge = make(map[string]float64)
var Storage forms.MemStorage
var M forms.Monitor
var Interval int = 1
var File string

type Producer struct {
	file *os.File

	writer *bufio.Writer
}

func ContainsMetrics(name string, types string) bool {
	var ok bool
	switch types {
	case "counter":
		_, ok = ServerCounter[name]
	case "gauge":
		_, ok = ServerGauge[name]
	}

	return ok

}

func NewProducer(filename string) (*Producer, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {

		return nil, fmt.Errorf("cannot open file: %v", err)
	}

	return &Producer{
		file: file,

		writer: bufio.NewWriter(file),
	}, nil
}

func (p *Producer) WriteEvent(metric *forms.Metrics) error {
	data, err := json.Marshal(&metric)
	if err != nil {

		return fmt.Errorf("cannot marshal: %v", err)
	}

	if _, err := p.writer.Write(data); err != nil {

		return fmt.Errorf("cannot write : %v", err)
	}

	if err := p.writer.WriteByte('\n'); err != nil {

		return fmt.Errorf("cannot write byte: %v", err)
	}

	return p.writer.Flush()
}

type Consumer struct {
	file *os.File

	reader *bufio.Reader
}

func NewConsumer(filename string) (*Consumer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)

	if err != nil {

		return nil, fmt.Errorf("cannot open file: %v", err)
	}

	return &Consumer{
		file: file,

		reader: bufio.NewReader(file),
	}, nil
}

func (c *Consumer) ReadEvent() (*forms.Metrics, error) {

	data, err := c.reader.ReadBytes('\n')

	if err != nil {

		return nil, fmt.Errorf("cannot read bytes: %v", err)
	}

	metric := forms.Metrics{}
	err = json.Unmarshal(data, &metric)
	if err != nil {

		return nil, fmt.Errorf("cannot unmarshal metrics: %v", err)
	}

	return &metric, nil
}
func (p *Producer) Close() error {
	return p.file.Close()
}
func (c *Consumer) Close() error {
	return c.file.Close()
}
func ConsInit(metric *forms.Metrics, fileName string) error {

	Producer, err := NewProducer(fileName)
	if err != nil {

		return fmt.Errorf("cannot create new produser: %v", err)
	}
	defer Producer.Close()

	if err := Producer.WriteEvent(metric); err != nil {

		return fmt.Errorf("cannot write event: %v", err)
	}
	return nil
}
func SaveMetrics(fileName string) error {
	var metrics forms.Metrics
	os.Remove(fileName)
	for key, value := range ServerGauge {

		metrics.Value = &value
		metrics.ID = key
		metrics.MType = "gauge"
		if err := ConsInit(&metrics, fileName); err != nil {

			return fmt.Errorf("cannot cons init: %v", err)
		}
	}
	for key, value := range ServerCounter {

		metrics.Delta = &value
		metrics.ID = key
		metrics.MType = "counter"
		if err := ConsInit(&metrics, fileName); err != nil {

			return fmt.Errorf("cannot cons init: %v", err)
		}
	}

	return nil
}
func SelectMetric(fileName string) error {

	defer os.Remove(fileName)
	var metric *forms.Metrics
	Producer, err := NewProducer(fileName)
	if err != nil {
		return fmt.Errorf("cannot create produser: %v", err)
	}
	defer Producer.Close()

	Consumer, err := NewConsumer(fileName)
	if err != nil {
		return fmt.Errorf("cannot create consumer: %v", err)
	}
	defer Consumer.Close()
	for {
		metric, err = Consumer.ReadEvent()

		if err != nil {

			return fmt.Errorf("cannot read event: %v", err)
		}
		metrics := *metric
		if metrics.MType == "counter" {
			ServerCounter[metric.ID] = *metric.Delta
		}
		if metrics.MType == "gauge" {
			ServerGauge[metric.ID] = *metric.Value
		}

	}

}

func Run(file string, interval int) {
	reportInterval := time.Duration(interval) * time.Second
	tick := time.NewTicker(reportInterval)
	if interval == 0 {
		Interval = interval
		File = file
		return
	}

	for {
		select {

		case <-tick.C:

			if err := SaveMetrics(file); err != nil {
				logger.LogError("cannot Save Metric: ", err)
			}
		case <-tick.C:

		}
	}
}
