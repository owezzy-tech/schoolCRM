package admissionsbus

import (
	"errors"
	"math"
	"regexp"
	"sort"
	"strings"
	"time"
)

var (
	ErrKenyaNationalIDInvalid   = errors.New("kenya national id invalid")
	ErrKenyaUPIInvalid          = errors.New("kenya upi invalid")
	ErrKenyaKCSEIndexInvalid    = errors.New("kenya kcse index number invalid")
	ErrKenyaPostalCodeInvalid   = errors.New("kenya postal code invalid")
	ErrKenyaAddressBlockMissing = errors.New("kenya address requires postal, physical, or geo block")
	ErrKCSEExamYearInvalid      = errors.New("kcse exam year invalid")
	ErrKCSESubjectCodeInvalid   = errors.New("kcse subject code invalid")
	ErrKCSEGradeInvalid         = errors.New("kcse grade invalid")
	ErrKCSESubjectsRequired     = errors.New("kcse subjects required")
	ErrKCSEMissingSubject       = errors.New("kcse missing required subject")
	ErrKCSEUngradedSubject      = errors.New("kcse subject is ungraded")
	ErrKCSEClusterInvalid       = errors.New("kcse cluster definition invalid")
)

var (
	kenyaNationalIDPattern = regexp.MustCompile(`^[0-9]{1,8}$`)
	kenyaUPIPattern        = regexp.MustCompile(`^[A-Z0-9]{6,12}$`)
	kenyaKCSEIndexPattern  = regexp.MustCompile(`^[0-9]{1,5}$`)
	kenyaPostalCodePattern = regexp.MustCompile(`^[0-9]{5}$`)
	kcseSubjectCodePattern = regexp.MustCompile(`^[0-9A-Z]{2,12}$`)
)

const (
	kcseMinExamYear         = 1989
	kcseBestSubjects        = 7
	kcseClusterSubjects     = 4
	kcseMaxAggregatePoints  = 84
	kcseMaxClusterRawPoints = 48
)

var kcseGradePoints = map[string]int{
	"A":  12,
	"A-": 11,
	"B+": 10,
	"B":  9,
	"B-": 8,
	"C+": 7,
	"C":  6,
	"C-": 5,
	"D+": 4,
	"D":  3,
	"D-": 2,
	"E":  1,
}

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

// KCSESubjectGrade represents a structurally valid KCSE subject grade.
type KCSESubjectGrade struct {
	SubjectCode string
	Grade       string
	Points      int
}

// ParseKCSESubjectGrade validates and normalizes a KCSE subject grade.
func ParseKCSESubjectGrade(subjectCode string, grade string) (KCSESubjectGrade, error) {
	normalizedSubjectCode := strings.ToUpper(strings.TrimSpace(subjectCode))
	if !kcseSubjectCodePattern.MatchString(normalizedSubjectCode) {
		return KCSESubjectGrade{}, ErrKCSESubjectCodeInvalid
	}

	normalizedGrade := strings.ToUpper(strings.TrimSpace(grade))
	points, exists := kcseGradePoints[normalizedGrade]
	if !exists {
		return KCSESubjectGrade{}, ErrKCSEGradeInvalid
	}

	return KCSESubjectGrade{
		SubjectCode: normalizedSubjectCode,
		Grade:       normalizedGrade,
		Points:      points,
	}, nil
}

// MustParseKCSESubjectGrade parses a KCSE subject grade and panics on failure.
func MustParseKCSESubjectGrade(subjectCode string, grade string) KCSESubjectGrade {
	subjectGrade, err := ParseKCSESubjectGrade(subjectCode, grade)
	if err != nil {
		panic(err)
	}

	return subjectGrade
}

// KCSEResult represents a structurally valid KCSE result used for admissions eligibility.
type KCSEResult struct {
	IndexNumber KenyaKCSEIndexNumber
	ExamYear    int
	Subjects    []KCSESubjectGrade
	MeanGrade   string
	MeanPoints  int
}

// ParseKCSEResult validates subject grades and calculates the official KNEC aggregate mean grade.
func ParseKCSEResult(indexNumber KenyaKCSEIndexNumber, examYear int, subjects []KCSESubjectGrade) (KCSEResult, error) {
	if examYear < kcseMinExamYear || examYear > time.Now().Year() {
		return KCSEResult{}, ErrKCSEExamYearInvalid
	}

	if len(subjects) == 0 {
		return KCSEResult{}, ErrKCSESubjectsRequired
	}

	normalizedSubjects := make([]KCSESubjectGrade, 0, len(subjects))
	for _, subject := range subjects {
		if !kcseSubjectCodePattern.MatchString(subject.SubjectCode) {
			return KCSEResult{}, ErrKCSESubjectCodeInvalid
		}

		if subject.Points <= 0 {
			return KCSEResult{}, ErrKCSEUngradedSubject
		}

		points, exists := kcseGradePoints[subject.Grade]
		if !exists || points != subject.Points {
			return KCSEResult{}, ErrKCSEGradeInvalid
		}

		normalizedSubjects = append(normalizedSubjects, subject)
	}

	aggregate := bestAggregatePoints(normalizedSubjects)
	meanGrade := kcseMeanGradeForAggregate(aggregate)

	return KCSEResult{
		IndexNumber: indexNumber,
		ExamYear:    examYear,
		Subjects:    normalizedSubjects,
		MeanGrade:   meanGrade,
		MeanPoints:  aggregate,
	}, nil
}

// KUCCPSClusterDefinition defines the four-subject cluster used for public cluster point approximation.
type KUCCPSClusterDefinition struct {
	Code         string
	Name         string
	SubjectCodes []string
}

// KCSEClusterWeight explains a deterministic public KUCCPS cluster-point approximation.
type KCSEClusterWeight struct {
	ClusterCode       string
	ClusterName       string
	RawClusterPoints  int
	AggregatePoints   int
	WeightedPoints    float64
	RequiredSubjects  []string
	ApproximationNote string
}

// ClusterWeight calculates the public KUCCPS cluster-point approximation for a result.
// KUCCPS states exact placement points rely on private KNEC performance indices, so this
// value is an auditable estimate based on the public sqrt((r/48)*(t/84))*48 formula.
func (result KCSEResult) ClusterWeight(cluster KUCCPSClusterDefinition) (KCSEClusterWeight, error) {
	if len(cluster.SubjectCodes) != kcseClusterSubjects {
		return KCSEClusterWeight{}, ErrKCSEClusterInvalid
	}

	gradeBySubject := make(map[string]KCSESubjectGrade, len(result.Subjects))
	for _, subject := range result.Subjects {
		gradeBySubject[subject.SubjectCode] = subject
	}

	rawClusterPoints := 0
	requiredSubjects := make([]string, 0, len(cluster.SubjectCodes))
	for _, subjectCode := range cluster.SubjectCodes {
		normalizedSubjectCode := strings.ToUpper(strings.TrimSpace(subjectCode))
		if !kcseSubjectCodePattern.MatchString(normalizedSubjectCode) {
			return KCSEClusterWeight{}, ErrKCSESubjectCodeInvalid
		}

		subject, exists := gradeBySubject[normalizedSubjectCode]
		if !exists {
			return KCSEClusterWeight{}, ErrKCSEMissingSubject
		}

		requiredSubjects = append(requiredSubjects, normalizedSubjectCode)
		rawClusterPoints += subject.Points
	}

	weightedPoints := math.Sqrt((float64(rawClusterPoints)/kcseMaxClusterRawPoints)*(float64(result.MeanPoints)/kcseMaxAggregatePoints)) * kcseMaxClusterRawPoints

	return KCSEClusterWeight{
		ClusterCode:       strings.TrimSpace(cluster.Code),
		ClusterName:       strings.TrimSpace(cluster.Name),
		RawClusterPoints:  rawClusterPoints,
		AggregatePoints:   result.MeanPoints,
		WeightedPoints:    roundToThree(weightedPoints),
		RequiredSubjects:  requiredSubjects,
		ApproximationNote: "Public KUCCPS approximation; exact KUCCPS points require private KNEC performance indices.",
	}, nil
}

func bestAggregatePoints(subjects []KCSESubjectGrade) int {
	points := make([]int, len(subjects))
	for i, subject := range subjects {
		points[i] = subject.Points
	}

	sort.Sort(sort.Reverse(sort.IntSlice(points)))

	limit := kcseBestSubjects
	if len(points) < limit {
		limit = len(points)
	}

	total := 0
	for i := 0; i < limit; i++ {
		total += points[i]
	}

	return total
}

func kcseMeanGradeForAggregate(aggregate int) string {
	switch {
	case aggregate >= 78:
		return "A"
	case aggregate >= 71:
		return "A-"
	case aggregate >= 64:
		return "B+"
	case aggregate >= 57:
		return "B"
	case aggregate >= 50:
		return "B-"
	case aggregate >= 43:
		return "C+"
	case aggregate >= 36:
		return "C"
	case aggregate >= 29:
		return "C-"
	case aggregate >= 22:
		return "D+"
	case aggregate >= 15:
		return "D"
	case aggregate >= 8:
		return "D-"
	default:
		return "E"
	}
}

func roundToThree(value float64) float64 {
	return math.Round(value*1000) / 1000
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
