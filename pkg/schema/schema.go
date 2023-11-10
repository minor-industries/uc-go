package schema

type SensorData struct {
	Temperature      float32 // celsius
	RelativeHumidity float32
}

type ThermocoupleData struct {
	Temperature float32 // celsius
	Description [16]byte
}
