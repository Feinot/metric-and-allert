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

		return nil, err
	}

	return &Producer{
		file: file,

		writer: bufio.NewWriter(file),
	}, nil
}

func (p *Producer) WriteEvent(metric *forms.Metrics) error {
	data, err := json.Marshal(&metric)
	if err != nil {

		return err
	}

	if _, err := p.writer.Write(data); err != nil {

		return err
	}

	if err := p.writer.WriteByte('\n'); err != nil {

		return err
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

		return nil, err
	}

	return &Consumer{
		file: file,

		reader: bufio.NewReader(file),
	}, nil
}

func (c *Consumer) ReadEvent() (*forms.Metrics, error) {

	data, err := c.reader.ReadBytes('\n')

	if err != nil {

		return nil, err
	}

	metric := forms.Metrics{}
	err = json.Unmarshal(data, &metric)
	if err != nil {

		return nil, err
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

		return (err)
	}
	defer Producer.Close()

	if err := Producer.WriteEvent(metric); err != nil {

		return (err)
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

			return err
		}
	}
	for key, value := range ServerCounter {

		metrics.Delta = &value
		metrics.ID = key
		metrics.MType = "counter"
		if err := ConsInit(&metrics, fileName); err != nil {

			return err
		}
	}

	return nil
}
func SelectMetric(fileName string) error {

	defer os.Remove(fileName)
	var metric *forms.Metrics
	Producer, err := NewProducer(fileName)
	if err != nil {
		return (err)
	}
	defer Producer.Close()

	Consumer, err := NewConsumer(fileName)
	if err != nil {
		return (err)
	}
	for q := 0; q < 1; {
		metric, err = Consumer.ReadEvent()

		if err != nil {
			q = 100
			return nil
		}
		metrics := *metric
		if metrics.MType == "counter" {
			ServerCounter[metric.ID] = *metric.Delta
		}
		if metrics.MType == "gauge" {
			ServerGauge[metric.ID] = *metric.Value
		}

	}
	defer Consumer.Close()

	return nil

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
			err := SaveMetrics(file)
			if err != nil {
				logger.LogError("cannot Save Metric: ", err)
			}
		case <-tick.C:
			fmt.Println(ServerGauge)
			fmt.Println(ServerCounter)
		}
	}
}
