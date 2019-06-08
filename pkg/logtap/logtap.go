package logtap

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/lichuan0620/logtap/pkg/fieldpath"
	"github.com/lichuan0620/logtap/pkg/logger"
)

// A LogTap is a runnable worker that keep generating log messages in a predefined way.
type LogTap interface {
	// GetTask is used to inspect the underlying task; it return a copy of the LogTask.
	GetTask() *LogTask

	// Run prompts the LogTap to start generating log messages and blocks until it stops. The LogTap would stop
	// when either the stopCh was closed or an error occurred. Run can only be called once per LogTap instance.
	Run(stopCh <-chan struct{}) error
}

type logTapImpl struct {
	task  *LogTask
	mutex sync.Mutex
	once  chan struct{}
}

// NewLogTap creates a LogTap with the given name; its behavior is defined by the given LogTaskSpec object.
func NewLogTap(taskTemplate *LogTaskSpec, name string) (LogTap, error) {
	if err := ValidateLogTaskSpec(fieldpath.NewFieldPath(), taskTemplate); err != nil {
		return nil, err
	}
	ret := &logTapImpl{
		task: &LogTask{
			Metadata: Metadata{
				Name:              name,
				CreationTimestamp: time.Now().UTC(),
			},
			Spec:   taskTemplate.DeepCopy(),
			Status: new(LogTaskStatus),
		},
		once: make(chan struct{}),
	}
	ret.setPhase(PhaseIdle, "")
	return ret, nil
}

func (lm logTapImpl) GetTask() *LogTask {
	lm.mutex.Lock()
	defer lm.mutex.Unlock()
	return lm.task.DeepCopy()
}

func (lm *logTapImpl) Run(stopCh <-chan struct{}) error {
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
		reason := fmt.Sprintf("[%s] unsupported output kind: %s", lm.task.Name, lm.task.Spec.OutputKind)
		lm.setPhase(PhaseFailed, reason)
		return fmt.Errorf(reason)
	}
	var worker logger.Logger
	switch lm.task.Spec.ContentType {
	case ContentTypeExplicit:
		worker = logger.NewExplicitLogger(output, lm.task.Spec.Message, lm.task.Name, lm.task.Spec.TimestampFormat)
	case ContentTypeRandom:
		worker = logger.NewRandomLogger(output, lm.task.Spec.MinSize, lm.task.Name, lm.task.Spec.TimestampFormat)
	default:
		reason := fmt.Sprintf("[%s] unsupported content type: %s", lm.task.Name, lm.task.Spec.ContentType)
		lm.setPhase(PhaseFailed, reason)
		return fmt.Errorf(reason)
	}
	lm.setPhase(PhaseRunning, "")
	interval := time.Duration(float64(time.Second) * lm.task.Spec.Interval)
	timer := time.NewTimer(0)
	defer timer.Stop()
	for {
		select {
		case <-stopCh:
			lm.setPhase(PhaseStopped, "")
			return nil
		case <-timer.C:
			timer.Reset(interval)
			_, size, err := worker.Log()
			if err != nil {
				reason := fmt.Sprintf("[%s] failed to write log: %s", lm.task.Name, err.Error())
				lm.setPhase(PhaseFailed, reason)
				return fmt.Errorf(reason)
			}
			lm.recordLogStatus(size)
		}
	}
}

func (lm *logTapImpl) recordLogStatus(size int) {
	lm.mutex.Lock()
	defer lm.mutex.Unlock()
	lm.task.Status.SentCount++
	lm.task.Status.SentBytes += int64(size)
}

func (lm *logTapImpl) setPhase(phase string, reason string) {
	lm.mutex.Lock()
	defer lm.mutex.Unlock()
	lm.task.Status.PhaseTimestamp = time.Now().UTC()
	lm.task.Status.Phase = phase
	lm.task.Status.Reason = reason
}
