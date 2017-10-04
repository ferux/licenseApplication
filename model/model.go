package model

//Model struct which contains an object that satisfy db interface
type Model struct {
	db
}
//New returns new Model
func New(db db) *Model {
	return &Model{db}
}

//SelectLicenses returns arrays of all licenses
func (m *Model) SelectLicenses() ([]*License, error) {return m.db.SelectLicenses()}
//SelectLicense return specified license
func (m *Model) SelectLicense(id string) (*License, error) {return m.db.SelectLicense(id)}
//UpdateLicense updates specified license
func (m *Model) UpdateLicense(l *License) (error) {return m.db.UpdateLicense(l)}