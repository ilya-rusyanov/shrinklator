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

	"github.com/ilya-rusyanov/shrinklator/internal/entities"

	"go.uber.org/zap"
)

// File - file storage
type File struct {
	data  map[string]string
	mutex sync.RWMutex
	log   Logger
	file  io.WriteCloser
}

type dto struct {
	ID    string `json:"uuid"`
	Short string `json:"short_url"`
	Long  string `json:"original_url"`
}

// NewFile constructs File object
func NewFile(log Logger, filename string) (*File, error) {
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

	return &File{
		data: initialValues,
		log:  log,
		file: dataSource,
	}, nil
}

// Put adds antry
func (f *File) Put(ctx context.Context, id, value string, uid *entities.UserID) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	err := f.memoryStore(id, value)
	if err != nil {
		return err
	}
	return f.fileStore(id, value)
}

// PutBatch adds multiple entries
func (f *File) PutBatch(ctx context.Context, data []entities.ShortLongPair) error {
	return fmt.Errorf("TODO")
}

// ByID searches entry by identifier
func (f *File) ByID(ctx context.Context, id string) (entities.ExpandResult, error) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	value, ok := f.data[id]

	if !ok {
		f.log.Info("cannot find record", zap.String("id", id))
		return entities.ExpandResult{}, ErrNotFound
	}

	f.log.Info("successuflly found record", zap.String("id", id),
		zap.String("value", value))

	return entities.ExpandResult{URL: value}, nil
}

// ByUID searches entry by user identifier
func (f *File) ByUID(context.Context, entities.UserID) (entities.PairArray, error) {
	return nil, errors.New("TODO")
}

// Delete deletes entry
func (f *File) Delete(context.Context, entities.DeleteRequest) error {
	return errors.New("TODO")
}

// MustClose finalizes storage
func (f *File) MustClose() {
	if err := f.file.Close(); err != nil {
		panic(fmt.Errorf("error closing file: %w", err))
	}
}

// Ping checks accessibility of storage
func (f *File) Ping(context.Context) error {
	return nil
}

// CountUsersAndUrls counts users and URLs
func (f *File) CountUsersAndUrls(context.Context) (entities.Stats, error) {
	return entities.Stats{}, errors.New("not implemented")
}

func (f *File) memoryStore(id, value string) error {
	if val, ok := f.data[id]; ok {
		return ErrAlreadyExists{
			StoredValue: val,
		}
	}

	f.data[id] = value

	return nil
}

func (f *File) fileStore(id, value string) error {
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

		obj := dto{}
		err = json.Unmarshal(data, &obj)
		if err != nil {
			return result, fmt.Errorf("error unmarshaling JSON: %w", err)
		}

		result[obj.Short] = obj.Long
		lastID++
	}

	return result, nil
}
