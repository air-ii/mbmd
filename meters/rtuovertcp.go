package meters

import (
	"time"

	"github.com/grid-x/modbus"
)

// RTUOverTCP is a RTU encoder over a TCP modbus connection
type RTUOverTCP struct {
	Client  modbus.Client
	Handler *modbus.RTUOverTCPClientHandler
	prevID  uint8
}

// NewRTUOverTCPClientHandler creates a RTU over TCP modbus handler
func NewRTUOverTCPClientHandler(device string) *modbus.RTUOverTCPClientHandler {
	handler := modbus.NewRTUOverTCPClientHandler(device)

	// set default timings
	handler.ProtocolRecoveryTimeout = 10 * time.Second // not used
	handler.LinkRecoveryTimeout = 15 * time.Second     // not used

	return handler
}

var _ Connection = (*RTUOverTCP)(nil)

// NewRTUOverTCP creates a TCP modbus client
func NewRTUOverTCP(address string) *RTUOverTCP {
	handler := NewRTUOverTCPClientHandler(address)
	client := modbus.NewClient(handler)

	b := &RTUOverTCP{
		Client:  client,
		Handler: handler,
	}

	return b
}

// String returns the bus connection address (TCP)
func (b *RTUOverTCP) String() string {
	return b.Handler.Address
}

// ModbusClient returns the TCP modbus client
func (b *RTUOverTCP) ModbusClient() modbus.Client {
	return b.Client
}

// Logger sets a logging instance for physical bus operations
func (b *RTUOverTCP) Logger(l Logger) {
	b.Handler.Logger = l
}

// Slave sets the modbus device id for the following operations
func (b *RTUOverTCP) Slave(deviceID uint8) {
	// Some devices like SDM need to have a little pause between querying different device ids
	if b.prevID != 0 && deviceID != b.prevID {
		time.Sleep(time.Duration(100) * time.Millisecond)
	}

	b.prevID = deviceID
	b.Handler.SetSlave(deviceID)
}

// Timeout sets the modbus timeout
func (b *RTUOverTCP) Timeout(timeout time.Duration) time.Duration {
	t := b.Handler.Timeout
	b.Handler.Timeout = timeout
	return t
}

// ConnectDelay sets the the initial delay after connecting before starting communication
func (b *RTUOverTCP) ConnectDelay(delay time.Duration) {
	b.Handler.ConnectDelay = delay
}

// Close closes the modbus connection.
// This forces the modbus client to reopen the connection before the next bus operations.
func (b *RTUOverTCP) Close() {
	b.Handler.Close()
}

// Clone clones the modbus connection.
func (b *RTUOverTCP) Clone(deviceID byte) Connection {
	handler := b.Handler.Clone()
	handler.SetSlave(deviceID)

	return &RTUOverTCP{
		Client:  modbus.NewClient(handler),
		Handler: handler,
	}
}
