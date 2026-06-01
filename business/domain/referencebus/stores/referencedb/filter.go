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
