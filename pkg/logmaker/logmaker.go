package logmaker

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/lichuan0620/logtap/pkg/fieldpath"
	"github.com/lichuan0620/logtap/pkg/logger"
)

// A LogMaker is a runnable worker that keep generating log messages in a predefined way.
type LogMaker interface {
	// GetTask is used to inspect the underlying task; it return a copy of the LogTask.
	GetTask() *LogTask

	// Run prompts the LogMaker to start generating logs and blocks until it stops. The LogMaker would stop when
	// either the stopCh was closed or an error occurred. Run can only be called once per LogMaker instance; a
	// second call would cause a panic.
	Run(stopCh <-chan struct{}) error
}

type logMakerImpl struct {
	task  *LogTask
	mutex sync.Mutex
	once  chan struct{}
}

// NewLogMaker creates a LogMaker with the given name; its behavior is defined by the given LogTaskSpec object.
func NewLogMaker(taskTemplate *LogTaskSpec, name string) (LogMaker, error) {
	if err := ValidateLogTaskSpec(fieldpath.NewFieldPath(), taskTemplate); err != nil {
		return nil, err
	}
	ret := &logMakerImpl{
		task: &LogTask{
			Metadata: Metadata{
				Name:              name,
				CreationTimestamp: time.Now(),
			},
			Spec:   taskTemplate.DeepCopy(),
			Status: new(LogTaskStatus),
		},
		once: make(chan struct{}),
	}
	return ret, nil
}

func (lm logMakerImpl) GetTask() *LogTask {
	lm.mutex.Lock()
	defer lm.mutex.Unlock()
	return lm.task.DeepCopy()
}

func (lm *logMakerImpl) Run(stopCh <-chan struct{}) error {
	close(lm.once)
	var output io.Writer
	switch lm.task.Spec.OutputKind {
	case OutputKindStdErr:
		output = os.Stderr
	case OutputKindStdOut:
		output = os.Stdout
	case OutputKindFile:
		file, err := os.OpenFile(lm.task.Spec.Filepath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("[%s] failed to open log file: %s", lm.task.Name, err.Error())
		}
		defer file.Close()
	default:
		return fmt.Errorf("[%s] unsupported output kind: %s", lm.task.Name, lm.task.Spec.OutputKind)
	}
	var worker logger.Logger
	switch lm.task.Spec.ContentType {
	case ContentTypeExplicit:
		worker = logger.NewExplicitLogger(output, lm.task.Spec.Message, lm.task.Name, lm.task.Spec.TimestampFormat)
	case ContentTypeRandom:
		worker = logger.NewRandomLogger(output, lm.task.Spec.MinSize, lm.task.Name, lm.task.Spec.TimestampFormat)
	default:
		return fmt.Errorf("[%s] unsupported content type: %s", lm.task.Name, lm.task.Spec.ContentType)
	}
	interval := time.Duration(float64(time.Second) * lm.task.Spec.Interval)
	timer := time.NewTimer(0)
	defer timer.Stop()
	for {
		select {
		case <-stopCh:
			return nil
		case <-timer.C:
			timer.Reset(interval)
			_, size, err := worker.Log()
			if err != nil {
				return fmt.Errorf("[%s] failed to write log: %s", lm.task.Name, err.Error())
			}
			lm.recordLogStatus(size)
		}
	}
}

func (lm *logMakerImpl) recordLogStatus(size int) {
	lm.mutex.Lock()
	defer lm.mutex.Unlock()
	lm.task.Status.SentCount++
	lm.task.Status.SentBytes += int64(size)
}