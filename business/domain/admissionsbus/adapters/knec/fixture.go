package knec

import (
	"time"

	"github.com/owezzy/schoolCRM/business/domain/admissionsbus"
)

// FixtureVerifyKCSEResultRequest is a representative KNEC verification request.
func FixtureVerifyKCSEResultRequest() VerifyKCSEResultRequest {
	return VerifyKCSEResultRequest{
		IndexNumber: admissionsbus.MustParseKenyaKCSEIndexNumber("12345"),
		ExamYear:    2024,
	}
}

// FixtureVerifyKCSEResultResponse is a representative KNEC verification response.
func FixtureVerifyKCSEResultResponse() VerifyKCSEResultResponse {
	result, err := admissionsbus.ParseKCSEResult(
		admissionsbus.MustParseKenyaKCSEIndexNumber("12345"),
		2024,
		[]admissionsbus.KCSESubjectGrade{
			admissionsbus.MustParseKCSESubjectGrade("ENG", "A"),
			admissionsbus.MustParseKCSESubjectGrade("KIS", "B+"),
			admissionsbus.MustParseKCSESubjectGrade("MAT", "A-"),
			admissionsbus.MustParseKCSESubjectGrade("BIO", "B"),
			admissionsbus.MustParseKCSESubjectGrade("CHE", "B+"),
			admissionsbus.MustParseKCSESubjectGrade("PHY", "B"),
			admissionsbus.MustParseKCSESubjectGrade("GEO", "B-"),
		},
	)
	if err != nil {
		panic(err)
	}

	return VerifyKCSEResultResponse{
		Verified:    true,
		Result:      result,
		Candidate:   "Amina Wanjiku",
		ExternalRef: "KNEC-VERIFY-2024-12345",
		VerifiedAt:  time.Date(2026, time.March, 3, 10, 0, 0, 0, time.UTC),
	}
}
