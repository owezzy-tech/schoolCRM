package admissionsbus

import (
	"errors"
	"math"
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

func TestParseKCSESubjectGrade(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		subjectCode string
		grade       string
		wantCode    string
		wantGrade   string
		wantPoints  int
		err         error
	}{
		{name: "valid subject grade", subjectCode: "101", grade: "A", wantCode: "101", wantGrade: "A", wantPoints: 12},
		{name: "normalizes subject and grade", subjectCode: " bio101 ", grade: " b+ ", wantCode: "BIO101", wantGrade: "B+", wantPoints: 10},
		{name: "invalid subject", subjectCode: "BIO-101", grade: "A", err: ErrKCSESubjectCodeInvalid},
		{name: "invalid grade", subjectCode: "101", grade: "X", err: ErrKCSEGradeInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := ParseKCSESubjectGrade(tt.subjectCode, tt.grade)
			if tt.err != nil {
				if !errors.Is(err, tt.err) {
					t.Fatalf("err = %v, want %v", err, tt.err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.SubjectCode != tt.wantCode || got.Grade != tt.wantGrade || got.Points != tt.wantPoints {
				t.Fatalf("got = %#v, want code=%q grade=%q points=%d", got, tt.wantCode, tt.wantGrade, tt.wantPoints)
			}
		})
	}
}

func TestParseKCSEResultCalculatesMeanGrade(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		subjects   []KCSESubjectGrade
		wantGrade  string
		wantPoints int
	}{
		{
			name: "best seven grade A boundary",
			subjects: []KCSESubjectGrade{
				MustParseKCSESubjectGrade("101", "A"),
				MustParseKCSESubjectGrade("102", "A"),
				MustParseKCSESubjectGrade("103", "A"),
				MustParseKCSESubjectGrade("104", "A"),
				MustParseKCSESubjectGrade("105", "A"),
				MustParseKCSESubjectGrade("106", "A"),
				MustParseKCSESubjectGrade("107", "A"),
				MustParseKCSESubjectGrade("108", "E"),
			},
			wantGrade:  "A",
			wantPoints: 84,
		},
		{
			name: "worked medicine example aggregate",
			subjects: []KCSESubjectGrade{
				MustParseKCSESubjectGrade("MAT", "B+"),
				MustParseKCSESubjectGrade("BIO", "A"),
				MustParseKCSESubjectGrade("CHE", "A-"),
				MustParseKCSESubjectGrade("ENG", "B"),
				MustParseKCSESubjectGrade("KIS", "B"),
				MustParseKCSESubjectGrade("GEO", "B-"),
				MustParseKCSESubjectGrade("HIS", "C+"),
			},
			wantGrade:  "B+",
			wantPoints: 66,
		},
		{
			name: "retake keeps best seven subjects",
			subjects: []KCSESubjectGrade{
				MustParseKCSESubjectGrade("101", "A"),
				MustParseKCSESubjectGrade("102", "A-"),
				MustParseKCSESubjectGrade("103", "B+"),
				MustParseKCSESubjectGrade("104", "B"),
				MustParseKCSESubjectGrade("105", "B-"),
				MustParseKCSESubjectGrade("106", "C+"),
				MustParseKCSESubjectGrade("107", "C"),
				MustParseKCSESubjectGrade("108", "E"),
				MustParseKCSESubjectGrade("109", "D"),
			},
			wantGrade:  "B",
			wantPoints: 63,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := ParseKCSEResult(MustParseKenyaKCSEIndexNumber("12345"), 2024, tt.subjects)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.MeanGrade != tt.wantGrade || got.MeanPoints != tt.wantPoints {
				t.Fatalf("got grade=%q points=%d, want grade=%q points=%d", got.MeanGrade, got.MeanPoints, tt.wantGrade, tt.wantPoints)
			}
		})
	}
}

func TestParseKCSEResultValidation(t *testing.T) {
	t.Parallel()

	validSubjects := []KCSESubjectGrade{
		MustParseKCSESubjectGrade("101", "A"),
		MustParseKCSESubjectGrade("102", "B+"),
	}

	tests := []struct {
		name     string
		year     int
		subjects []KCSESubjectGrade
		want     error
	}{
		{name: "invalid year", year: 1988, subjects: validSubjects, want: ErrKCSEExamYearInvalid},
		{name: "subjects required", year: 2024, want: ErrKCSESubjectsRequired},
		{name: "ungraded subject", year: 2024, subjects: []KCSESubjectGrade{{SubjectCode: "101", Grade: "A"}}, want: ErrKCSEUngradedSubject},
		{name: "invalid grade points", year: 2024, subjects: []KCSESubjectGrade{{SubjectCode: "101", Grade: "A", Points: 11}}, want: ErrKCSEGradeInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := ParseKCSEResult(MustParseKenyaKCSEIndexNumber("12345"), tt.year, tt.subjects)
			if !errors.Is(err, tt.want) {
				t.Fatalf("err = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestKCSEResultClusterWeight(t *testing.T) {
	t.Parallel()

	result, err := ParseKCSEResult(MustParseKenyaKCSEIndexNumber("12345"), 2024, []KCSESubjectGrade{
		MustParseKCSESubjectGrade("MAT", "B+"),
		MustParseKCSESubjectGrade("BIO", "A"),
		MustParseKCSESubjectGrade("CHE", "A-"),
		MustParseKCSESubjectGrade("ENG", "B"),
		MustParseKCSESubjectGrade("KIS", "B"),
		MustParseKCSESubjectGrade("GEO", "B-"),
		MustParseKCSESubjectGrade("HIS", "C+"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cluster := KUCCPSClusterDefinition{
		Code:         "13",
		Name:         "Medicine and Health Sciences",
		SubjectCodes: []string{"MAT", "BIO", "CHE", "ENG"},
	}

	got, err := result.ClusterWeight(cluster)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.RawClusterPoints != 42 {
		t.Fatalf("raw cluster points = %d, want 42", got.RawClusterPoints)
	}
	if got.AggregatePoints != 66 {
		t.Fatalf("aggregate points = %d, want 66", got.AggregatePoints)
	}
	if math.Abs(got.WeightedPoints-39.799) > 0.001 {
		t.Fatalf("weighted points = %.3f, want 39.799", got.WeightedPoints)
	}
	if got.ApproximationNote == "" {
		t.Fatal("approximation note is empty")
	}
}

func TestKCSEResultClusterWeightValidation(t *testing.T) {
	t.Parallel()

	result, err := ParseKCSEResult(MustParseKenyaKCSEIndexNumber("12345"), 2024, []KCSESubjectGrade{
		MustParseKCSESubjectGrade("MAT", "B+"),
		MustParseKCSESubjectGrade("BIO", "A"),
		MustParseKCSESubjectGrade("CHE", "A-"),
		MustParseKCSESubjectGrade("ENG", "B"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tests := []struct {
		name    string
		cluster KUCCPSClusterDefinition
		want    error
	}{
		{name: "invalid cluster size", cluster: KUCCPSClusterDefinition{SubjectCodes: []string{"MAT"}}, want: ErrKCSEClusterInvalid},
		{name: "invalid subject code", cluster: KUCCPSClusterDefinition{SubjectCodes: []string{"MAT", "BIO", "CHE", "ENG-1"}}, want: ErrKCSESubjectCodeInvalid},
		{name: "missing subject", cluster: KUCCPSClusterDefinition{SubjectCodes: []string{"MAT", "BIO", "CHE", "PHY"}}, want: ErrKCSEMissingSubject},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := result.ClusterWeight(tt.cluster)
			if !errors.Is(err, tt.want) {
				t.Fatalf("err = %v, want %v", err, tt.want)
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
