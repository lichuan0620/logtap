package v1alpha1

import "time"

// Version defines the version of the models in this package.
const Version = "v1alpha1"

// Metadata stores the metadata of a LogTask.
type Metadata struct {
	Version           string    `json:"version"`
	Name              string    `json:"name"`
	CreationTimestamp time.Time `json:"creationTimestamp,omitempty"`
}
