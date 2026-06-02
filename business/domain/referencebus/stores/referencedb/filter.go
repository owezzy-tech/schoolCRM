package referencedb

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/owezzy/schoolCRM/business/domain/referencebus"
)

func (s *Store) applyCountyFilter(filter referencebus.QueryFilter, data map[string]any, buf *bytes.Buffer) {
	var wc []string

	if filter.Code != nil {
		data["code"] = filter.Code
		wc = append(wc, "code = :code")
	}

	if filter.Name != nil {
		data["name"] = fmt.Sprintf("%%%s%%", *filter.Name)
		wc = append(wc, "name LIKE :name")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}

func (s *Store) applySubCountyFilter(filter referencebus.QueryFilter, data map[string]any, buf *bytes.Buffer) {
	var wc []string

	if filter.Code != nil {
		data["code"] = filter.Code
		wc = append(wc, "code = :code")
	}

	if filter.CountyCode != nil {
		data["county_code"] = filter.CountyCode
		wc = append(wc, "county_code = :county_code")
	}

	if filter.Name != nil {
		data["name"] = fmt.Sprintf("%%%s%%", *filter.Name)
		wc = append(wc, "name LIKE :name")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}

func (s *Store) applyWardFilter(filter referencebus.QueryFilter, data map[string]any, buf *bytes.Buffer) {
	var wc []string

	if filter.Code != nil {
		data["code"] = filter.Code
		wc = append(wc, "code = :code")
	}

	if filter.CountyCode != nil {
		data["county_code"] = filter.CountyCode
		wc = append(wc, "county_code = :county_code")
	}

	if filter.SubCountyCode != nil {
		data["sub_county_code"] = filter.SubCountyCode
		wc = append(wc, "sub_county_code = :sub_county_code")
	}

	if filter.Name != nil {
		data["name"] = fmt.Sprintf("%%%s%%", *filter.Name)
		wc = append(wc, "name LIKE :name")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}

func (s *Store) applyUniversityFilter(filter referencebus.QueryFilter, data map[string]any, buf *bytes.Buffer) {
	var wc []string

	if filter.Code != nil {
		data["code"] = filter.Code
		wc = append(wc, "code = :code")
	}

	if filter.Name != nil {
		data["name"] = fmt.Sprintf("%%%s%%", *filter.Name)
		wc = append(wc, "name LIKE :name")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}

func (s *Store) applyClusterFilter(filter referencebus.QueryFilter, data map[string]any, buf *bytes.Buffer) {
	var wc []string

	if filter.Code != nil {
		data["code"] = filter.Code
		wc = append(wc, "code = :code")
	}

	if filter.Name != nil {
		data["name"] = fmt.Sprintf("%%%s%%", *filter.Name)
		wc = append(wc, "name LIKE :name")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}

func (s *Store) applyKNQFLevelFilter(filter referencebus.QueryFilter, data map[string]any, buf *bytes.Buffer) {
	var wc []string

	if filter.Code != nil {
		data["code"] = filter.Code
		wc = append(wc, "code = :code")
	}

	if filter.Name != nil {
		data["name"] = fmt.Sprintf("%%%s%%", *filter.Name)
		wc = append(wc, "name LIKE :name")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}

func (s *Store) applyProgrammeFilter(filter referencebus.QueryFilter, data map[string]any, buf *bytes.Buffer) {
	var wc []string

	if filter.Code != nil {
		data["code"] = filter.Code
		wc = append(wc, "code = :code")
	}

	if filter.UniversityCode != nil {
		data["university_code"] = filter.UniversityCode
		wc = append(wc, "university_code = :university_code")
	}

	if filter.ClusterCode != nil {
		data["cluster_code"] = filter.ClusterCode
		wc = append(wc, "cluster_code = :cluster_code")
	}

	if filter.KNQFLevelCode != nil {
		data["knqf_level_code"] = filter.KNQFLevelCode
		wc = append(wc, "knqf_level_code = :knqf_level_code")
	}

	if filter.Name != nil {
		data["name"] = fmt.Sprintf("%%%s%%", *filter.Name)
		wc = append(wc, "name LIKE :name")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}
