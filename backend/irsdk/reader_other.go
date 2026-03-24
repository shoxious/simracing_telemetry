//go:build !windows

package irsdk

import "errors"

// errNotWindows is returned when running outside Windows.
var errNotWindows = errors.New("iRacing SDK is only available on Windows")

type stubReader struct{}

// NewWindowsReader returns a non-functional stub on non-Windows systems.
// Use --simulate flag instead.
func NewWindowsReader() Reader {
	return &stubReader{}
}

func (s *stubReader) Connect() error             { return errNotWindows }
func (s *stubReader) IsConnected() bool          { return false }
func (s *stubReader) ReadFrame() (*TelemetryFrame, error) { return nil, errNotWindows }
func (s *stubReader) SessionYAML() (string, error) { return "", errNotWindows }
func (s *stubReader) SessionUpdateCount() int32  { return 0 }
func (s *stubReader) Close()                     {}
