package model

type db interface{
	SelectLicenses() ([]*License, error)
	SelectLicense(string) (*License, error)
	UpdateLicense(*License) error
}

