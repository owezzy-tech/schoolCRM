package admissionsbus

import (
	"errors"
	"testing"
)

func TestParseKenyaNationalID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value string
		want  string
		err   error
	}{
		{name: "one digit", value: "1", want: "1"},
		{name: "eight digits", value: "12345678", want: "12345678"},
		{name: "trims whitespace", value: " 1234567 ", want: "1234567"},
		{name: "empty", value: "", err: ErrKenyaNationalIDInvalid},
		{name: "too long", value: "123456789", err: ErrKenyaNationalIDInvalid},
		{name: "letters", value: "A1234567", err: ErrKenyaNationalIDInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := ParseKenyaNationalID(tt.value)
			if tt.err != nil {
				if !errors.Is(err, tt.err) {
					t.Fatalf("err = %v, want %v", err, tt.err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.String() != tt.want {
				t.Fatalf("got = %q, want %q", got.String(), tt.want)
			}
		})
	}
}

func TestParseKenyaUPI(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value string
		want  string
		err   error
	}{
		{name: "six alphanumeric", value: "ABC123", want: "ABC123"},
		{name: "normalizes lowercase", value: "upi12345", want: "UPI12345"},
		{name: "twelve characters", value: "A1B2C3D4E5F6", want: "A1B2C3D4E5F6"},
		{name: "too short", value: "A1234", err: ErrKenyaUPIInvalid},
		{name: "too long", value: "A1B2C3D4E5F6G", err: ErrKenyaUPIInvalid},
		{name: "symbol", value: "UPI-123", err: ErrKenyaUPIInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := ParseKenyaUPI(tt.value)
			if tt.err != nil {
				if !errors.Is(err, tt.err) {
					t.Fatalf("err = %v, want %v", err, tt.err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.String() != tt.want {
				t.Fatalf("got = %q, want %q", got.String(), tt.want)
			}
		})
	}
}

func TestParseKenyaKCSEIndexNumber(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value string
		want  string
		err   error
	}{
		{name: "one digit", value: "7", want: "7"},
		{name: "five digits", value: "12345", want: "12345"},
		{name: "trims whitespace", value: " 2345 ", want: "2345"},
		{name: "empty", value: "", err: ErrKenyaKCSEIndexInvalid},
		{name: "too long", value: "123456", err: ErrKenyaKCSEIndexInvalid},
		{name: "letters", value: "12A45", err: ErrKenyaKCSEIndexInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := ParseKenyaKCSEIndexNumber(tt.value)
			if tt.err != nil {
				if !errors.Is(err, tt.err) {
					t.Fatalf("err = %v, want %v", err, tt.err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.String() != tt.want {
				t.Fatalf("got = %q, want %q", got.String(), tt.want)
			}
		})
	}
}

func TestKenyaAddressValidate(t *testing.T) {
	t.Parallel()

	poBox := "12345"
	postalCode := "00100"
	invalidPostalCode := "100"
	postalTown := "Nairobi"
	countyName := "Nairobi"
	ward := "Kilimani"
	latitude := -1.2921
	longitude := 36.8219
	plusCode := "6GCRPR5C+5X"

	tests := []struct {
		name    string
		address KenyaAddress
		want    error
	}{
		{
			name: "postal block",
			address: KenyaAddress{
				POBox:      &poBox,
				PostalCode: &postalCode,
				PostalTown: &postalTown,
			},
		},
		{
			name: "physical block",
			address: KenyaAddress{
				CountyName: &countyName,
				Ward:       &ward,
			},
		},
		{
			name: "geo coordinates block",
			address: KenyaAddress{
				Latitude:  &latitude,
				Longitude: &longitude,
			},
		},
		{
			name: "plus code block",
			address: KenyaAddress{
				PlusCode: &plusCode,
			},
		},
		{
			name: "invalid postal code",
			address: KenyaAddress{
				PostalCode: &invalidPostalCode,
			},
			want: ErrKenyaPostalCodeInvalid,
		},
		{
			name:    "missing all blocks",
			address: KenyaAddress{},
			want:    ErrKenyaAddressBlockMissing,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.address.Validate()
			if tt.want != nil {
				if !errors.Is(err, tt.want) {
					t.Fatalf("err = %v, want %v", err, tt.want)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestKenyaValueObjectsMarshalText(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		text []byte
	}{
		{name: "national id", text: mustMarshalText(t, MustParseKenyaNationalID("12345678"))},
		{name: "upi", text: mustMarshalText(t, MustParseKenyaUPI("abc12345"))},
		{name: "kcse index", text: mustMarshalText(t, MustParseKenyaKCSEIndexNumber("12345"))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if len(tt.text) == 0 {
				t.Fatal("text is empty")
			}
		})
	}
}

type textMarshaler interface {
	MarshalText() ([]byte, error)
}

func mustMarshalText(t *testing.T, value textMarshaler) []byte {
	t.Helper()

	text, err := value.MarshalText()
	if err != nil {
		t.Fatalf("MarshalText returned error: %v", err)
	}

	return text
}
