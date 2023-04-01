package interfaces

type ISerializable interface {
	Serialize() ([]byte, error)
	Deserialize(data []byte) error
}
