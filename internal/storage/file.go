package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"

	"github.com/ilya-rusyanov/shrinklator/internal/logger"

	"go.uber.org/zap"
)

type file struct {
	data  map[string]string
	mutex sync.RWMutex
	log   *logger.Log
	file  io.WriteCloser
}

type dto struct {
	ID    string `json:"uuid"`
	Short string `json:"short_url"`
	Long  string `json:"original_url"`
}

func NewFile(log *logger.Log, filename string) (*file, error) {
	dataSource, err := os.OpenFile(filename,
		os.O_RDWR|os.O_APPEND|os.O_CREATE,
		0640)

	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}

	log.Info("opened file persistence",
		zap.String("file path", filename))

	initialValues, err := readAll(dataSource)
	if err != nil {
		return nil, fmt.Errorf("error reading values: %w", err)
	}

	return &file{
		data: initialValues,
		log:  log,
		file: dataSource,
	}, nil
}

func (f *file) Put(ctx context.Context, id, value string) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.memoryStore(id, value)
	return f.fileStore(id, value)
}

func (f *file) ByID(ctx context.Context, id string) (string, error) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	value, ok := f.data[id]

	if !ok {
		f.log.Info("cannot find record", zap.String("id", id))
		return "", ErrNotFound
	}

	f.log.Info("successuflly found record", zap.String("id", id),
		zap.String("value", value))

	return value, nil
}

func (f *file) Close() {
	f.file.Close()
}

func (f *file) memoryStore(id, value string) {
	f.data[id] = value
}

func (f *file) fileStore(id, value string) error {
	writer := bufio.NewWriter(f.file)

	payload := dto{Short: id, Long: value, ID: strconv.Itoa(len(f.data))}
	encoded, err := json.Marshal(&payload)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	if _, err := writer.Write(encoded); err != nil {
		return fmt.Errorf("error writing encoded data: %w", err)
	}

	if err := writer.WriteByte('\n'); err != nil {
		return fmt.Errorf("error adding newline: %w", err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("error flushing file: %w", err)
	}

	return nil
}

func readAll(from io.Reader) (map[string]string, error) {
	result := make(map[string]string)

	lastID := 0
	reader := bufio.NewReader(from)
	data, err := reader.ReadBytes('\n')
	for ; !errors.Is(err, io.EOF); data, err = reader.ReadBytes('\n') {
		if err != nil {
			return result, fmt.Errorf("error reading line: %w", err)
		}

		dto := dto{}
		err = json.Unmarshal(data, &dto)
		if err != nil {
			return result, fmt.Errorf("error unmarshaling JSON: %w", err)
		}

		result[dto.Short] = dto.Long
		lastID++
	}

	return result, nil
}
