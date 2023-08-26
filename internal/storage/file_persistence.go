package storage

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/ilya-rusyanov/shrinklator/internal/logger"
	"go.uber.org/zap"
)

type FilePersistence struct {
	file   *os.File
	writer *bufio.Writer
	lastID int
}

type dataPiece struct {
	id    int    `json:"uuid"`
	short string `json:"short_url"`
	long  string `json:"original_url"`
}

func NewFilePersistence(filePath string) (*FilePersistence, error) {
	file, err := os.OpenFile(filePath,
		os.O_RDWR|os.O_APPEND|os.O_CREATE,
		0640)

	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}

	res := &FilePersistence{
		file:   file,
		writer: bufio.NewWriter(file),
	}

	logger.Log.Info("opened file persistence",
		zap.String("file path", filePath))
	return res, nil
}

func (p *FilePersistence) Append(short string, long string) error {
	var err error
	p.lastID++

	piece := dataPiece{short: short, long: long}
	encoded, err := json.Marshal(&piece)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	if _, err := p.writer.Write(encoded); err != nil {
		return fmt.Errorf("error writing encoded data: %w", err)
	}

	if err := p.writer.WriteByte('\n'); err != nil {
		return fmt.Errorf("error adding newline: %w", err)
	}

	if err := p.writer.Flush(); err != nil {
		return fmt.Errorf("error flushing file: %w", err)
	}

	return nil
}

func (p *FilePersistence) ReadAll() (values map[string]string, err error) {
	values = make(map[string]string)

	reader := bufio.NewReader(p.file)
	data := make([]byte, 0)
	for data, err = reader.ReadBytes('\n'); !errors.Is(err, io.EOF); data, err = reader.ReadBytes('\n') {
		if err != nil {
			err = fmt.Errorf("error reading line: %w", err)
			return
		}

		piece := dataPiece{}
		err = json.Unmarshal(data, &piece)
		if err != nil {
			err = fmt.Errorf("error unmarshaling JSON: %w", err)
			return
		}

		values[piece.short] = piece.long
		p.lastID++
	}

	// clear EOF error
	err = nil

	return
}

func (p *FilePersistence) Close() {
	p.file.Close()
}
