package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Feinot/metric-and-allert/internal/forms"
)

var ServerCounter = make(map[string]int64)
var ServerGauge = make(map[string]float64)
var AgentCounter = make(map[string]int64)
var AgentGauge = make(map[string]float64)
var Storage forms.MemStorage
var M forms.Monitor

type Producer struct {
	file *os.File
	// добавляем Writer в Producer
	writer *bufio.Writer
}

func NewProducer(filename string) (*Producer, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &Producer{
		file: file,
		// создаём новый Writer
		writer: bufio.NewWriter(file),
	}, nil
}

func (p *Producer) WriteEvent(metric *forms.Metrics) error {
	data, err := json.Marshal(&metric)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// записываем событие в буфер
	if _, err := p.writer.Write(data); err != nil {
		fmt.Println(err)
		return err
	}

	// добавляем перенос строки
	if err := p.writer.WriteByte('\n'); err != nil {
		fmt.Println(err)
		return err
	}

	// записываем буфер в файл
	return p.writer.Flush()
}

type Consumer struct {
	file *os.File
	// добавляем reader в Consumer
	reader *bufio.Reader
}

func NewConsumer(filename string) (*Consumer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &Consumer{
		file: file,
		// создаём новый Reader
		reader: bufio.NewReader(file),
	}, nil
}

func (c *Consumer) ReadEvent() (*forms.Metrics, error) {
	// читаем данные до символа переноса строки
	data, err := c.reader.ReadBytes('\n')

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// преобразуем данные из JSON-представления в структуру
	metric := forms.Metrics{}
	err = json.Unmarshal(data, &metric)
	if err != nil {
		fmt.Println(metric)
		return nil, err
	}
	fmt.Println(metric)
	return &metric, nil
}
func (p *Producer) Close() error {
	return p.file.Close()
}
func (c *Consumer) Close() error {
	return c.file.Close()
}
func ConsInit(metric *forms.Metrics, fileName string) error {

	//defer os.Remove(fileName)

	Producer, err := NewProducer(fileName)
	if err != nil {
		fmt.Println(err)
		return (err)
	}
	defer Producer.Close()

	if err := Producer.WriteEvent(metric); err != nil {
		fmt.Println(err)
		return (err)
	}
	return nil
}
func SaveMetrics(fileName string) error {
	var metrics forms.Metrics
	for key, value := range ServerGauge {

		metrics.Value = &value
		metrics.ID = key
		metrics.MType = "gauge"
		if err := ConsInit(&metrics, fileName); err != nil {
			fmt.Println(err)
			return err
		}
	}
	for key, value := range ServerCounter {

		metrics.Delta = &value
		metrics.ID = key
		metrics.MType = "caunter"
		if err := ConsInit(&metrics, fileName); err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}
func SelectMetric(fileName string) error {

	//defer os.Remove(fileName)
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
		if metrics.MType == "caunter" {
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

	for {
		select {

		case <-tick.C:
			fmt.Println(SaveMetrics(file))

		}
	}
}
