package logmaker

import (
	"fmt"
	"github.com/lichuan0620/logtap/pkg/fieldpath"
)

// ValidateLogTask validates a LogTask object.
func ValidateLogTask(path fieldpath.FieldPath, task *LogTask) (err error) {
	if err = ValidateLogTaskSpec(path.Add("spec"), task.Spec); err != nil {
		return err
	}
	if err = ValidateLogTaskStatus(path.Add("status"), task.Status); err != nil {
		return err
	}
	return nil
}

// ValidateLogTaskSpec validates a LogTaskSpec object.
func ValidateLogTaskSpec(path fieldpath.FieldPath, spec *LogTaskSpec) error {
	switch spec.ContentType {
	case ContentTypeRandom:
		if len(spec.Message) > 0 {
			return newValidationError(path.Add("message").String(), "invalid field")
		}
		if spec.MinSize < 0 {
			return newInvalidValueError(path.Add("minSize").String())
		}
	case ContentTypeExplicit:
		if spec.MinSize != 0 {
			return newValidationError(path.Add("minSize").String(), "invalid field")
		}
	default:
		return newValidationError(path.Add("contentType").String(), "unrecognized contentType")
	}
	filepathProvided := len(spec.Filepath) > 0
	switch spec.OutputKind {
	case OutputKindFile:
		if !filepathProvided {
			return newValidationError(path.Add(
				"filepath").String(),
				"filepath not specified for file output",
			)
		}
	case OutputKindStdErr:
		if filepathProvided {
			return newValidationError(path.Add(
				"filepath").String(),
				"filepath specified for STDERR output",
			)
		}
	case OutputKindStdOut:
		if filepathProvided {
			return newValidationError(path.Add(
				"filepath").String(),
				"filepath specified for STDOUT output",
			)
		}
	default:
		return newValidationError(path.Add("outputKind").String(), "unrecognized output kind")
	}
	if spec.Interval < 0 {
		return newInvalidValueError(path.Add("interval").String())
	}
	return nil
}

// ValidateLogTaskStatus validates a LogTaskStatus object.
func ValidateLogTaskStatus(path fieldpath.FieldPath, status *LogTaskStatus) error {
	if status.SentCount < 0 {
		return newInvalidValueError(path.Add("sentCount").String())
	}
	if status.SentBytes < 0 {
		return newInvalidValueError(path.Add("sentBytes").String())
	}
	return nil
}

func newValidationError(path, reason string) error {
	return fmt.Errorf("invalid field <%s>: %s", path, reason)
}

func newInvalidValueError(path string) error {
	return newValidationError(path, "invalid value")
}
