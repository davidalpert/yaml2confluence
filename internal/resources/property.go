package resources

type Property struct {
	id      string
	key     string
	value   string
	version int
}

func NewProperty(id, key, value string, version int) Property {
	return Property{id, key, value, version}
}
func (p Property) GetId() string {
	return p.id
}
func (p Property) GetKey() string {
	return p.key
}
func (p Property) GetValue() string {
	return p.value
}
func (p Property) GetIncrementedVersion() int {
	return p.version + 1
}
func (p Property) IsUpdate() bool {
	return p.id != "" && p.version > 0
}
