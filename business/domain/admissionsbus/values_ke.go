package admissionsbus

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrKenyaNationalIDInvalid   = errors.New("kenya national id invalid")
	ErrKenyaUPIInvalid          = errors.New("kenya upi invalid")
	ErrKenyaKCSEIndexInvalid    = errors.New("kenya kcse index number invalid")
	ErrKenyaPostalCodeInvalid   = errors.New("kenya postal code invalid")
	ErrKenyaAddressBlockMissing = errors.New("kenya address requires postal, physical, or geo block")
)

var (
	kenyaNationalIDPattern = regexp.MustCompile(`^[0-9]{1,8}$`)
	kenyaUPIPattern        = regexp.MustCompile(`^[A-Z0-9]{6,12}$`)
	kenyaKCSEIndexPattern  = regexp.MustCompile(`^[0-9]{1,5}$`)
	kenyaPostalCodePattern = regexp.MustCompile(`^[0-9]{5}$`)
)

// KenyaNationalID represents a structurally valid Kenyan national identity number.
type KenyaNationalID struct {
	value string
}

// ParseKenyaNationalID validates and normalizes a Kenyan national identity number.
func ParseKenyaNationalID(value string) (KenyaNationalID, error) {
	normalized := strings.TrimSpace(value)
	if !kenyaNationalIDPattern.MatchString(normalized) {
		return KenyaNationalID{}, ErrKenyaNationalIDInvalid
	}

	return KenyaNationalID{value: normalized}, nil
}

// MustParseKenyaNationalID parses a Kenyan national identity number and panics on failure.
func MustParseKenyaNationalID(value string) KenyaNationalID {
	id, err := ParseKenyaNationalID(value)
	if err != nil {
		panic(err)
	}

	return id
}

// String returns the normalized national identity number.
func (id KenyaNationalID) String() string {
	return id.value
}

// Equal reports whether two national identity numbers are equal.
func (id KenyaNationalID) Equal(other KenyaNationalID) bool {
	return id.value == other.value
}

// MarshalText implements encoding.TextMarshaler.
func (id KenyaNationalID) MarshalText() ([]byte, error) {
	return []byte(id.value), nil
}

// KenyaUPI represents a structurally valid NEMIS/KEMIS learner unique personal identifier.
type KenyaUPI struct {
	value string
}

// ParseKenyaUPI validates and normalizes a Kenyan learner UPI.
func ParseKenyaUPI(value string) (KenyaUPI, error) {
	normalized := strings.ToUpper(strings.TrimSpace(value))
	if !kenyaUPIPattern.MatchString(normalized) {
		return KenyaUPI{}, ErrKenyaUPIInvalid
	}

	return KenyaUPI{value: normalized}, nil
}

// MustParseKenyaUPI parses a learner UPI and panics on failure.
func MustParseKenyaUPI(value string) KenyaUPI {
	upi, err := ParseKenyaUPI(value)
	if err != nil {
		panic(err)
	}

	return upi
}

// String returns the normalized learner UPI.
func (upi KenyaUPI) String() string {
	return upi.value
}

// Equal reports whether two learner UPIs are equal.
func (upi KenyaUPI) Equal(other KenyaUPI) bool {
	return upi.value == other.value
}

// MarshalText implements encoding.TextMarshaler.
func (upi KenyaUPI) MarshalText() ([]byte, error) {
	return []byte(upi.value), nil
}

// KenyaKCSEIndexNumber represents a structurally valid KCSE candidate index number.
type KenyaKCSEIndexNumber struct {
	value string
}

// ParseKenyaKCSEIndexNumber validates and normalizes a KCSE index number.
func ParseKenyaKCSEIndexNumber(value string) (KenyaKCSEIndexNumber, error) {
	normalized := strings.TrimSpace(value)
	if !kenyaKCSEIndexPattern.MatchString(normalized) {
		return KenyaKCSEIndexNumber{}, ErrKenyaKCSEIndexInvalid
	}

	return KenyaKCSEIndexNumber{value: normalized}, nil
}

// MustParseKenyaKCSEIndexNumber parses a KCSE index number and panics on failure.
func MustParseKenyaKCSEIndexNumber(value string) KenyaKCSEIndexNumber {
	index, err := ParseKenyaKCSEIndexNumber(value)
	if err != nil {
		panic(err)
	}

	return index
}

// String returns the normalized KCSE index number.
func (index KenyaKCSEIndexNumber) String() string {
	return index.value
}

// Equal reports whether two KCSE index numbers are equal.
func (index KenyaKCSEIndexNumber) Equal(other KenyaKCSEIndexNumber) bool {
	return index.value == other.value
}

// MarshalText implements encoding.TextMarshaler.
func (index KenyaKCSEIndexNumber) MarshalText() ([]byte, error) {
	return []byte(index.value), nil
}

// KenyaAddress represents a Kenyan postal, physical, or geocoded address.
type KenyaAddress struct {
	POBox            *string
	PostalCode       *string
	PostalTown       *string
	CountyCode       *string
	CountyName       *string
	SubCounty        *string
	Ward             *string
	EstateOrLocality *string
	Street           *string
	Building         *string
	HouseNumber      *string
	Latitude         *float64
	Longitude        *float64
	PlusCode         *string
}

// Validate enforces the Kenya address shape that can be checked without reference data.
func (address KenyaAddress) Validate() error {
	if address.PostalCode != nil {
		postalCode := strings.TrimSpace(*address.PostalCode)
		if !kenyaPostalCodePattern.MatchString(postalCode) {
			return ErrKenyaPostalCodeInvalid
		}
	}

	if !address.hasPostalBlock() && !address.hasPhysicalBlock() && !address.hasGeoBlock() {
		return ErrKenyaAddressBlockMissing
	}

	return nil
}

func (address KenyaAddress) hasPostalBlock() bool {
	return hasString(address.POBox) || hasString(address.PostalCode) || hasString(address.PostalTown)
}

func (address KenyaAddress) hasPhysicalBlock() bool {
	return hasString(address.CountyCode) ||
		hasString(address.CountyName) ||
		hasString(address.SubCounty) ||
		hasString(address.Ward) ||
		hasString(address.EstateOrLocality) ||
		hasString(address.Street) ||
		hasString(address.Building) ||
		hasString(address.HouseNumber)
}

func (address KenyaAddress) hasGeoBlock() bool {
	return address.Latitude != nil || address.Longitude != nil || hasString(address.PlusCode)
}

func hasString(value *string) bool {
	return value != nil && strings.TrimSpace(*value) != ""
}
