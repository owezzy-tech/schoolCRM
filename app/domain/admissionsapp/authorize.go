package admissionsapp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/owezzy/schoolCRM/app/sdk/authclient"
	"github.com/owezzy/schoolCRM/app/sdk/errs"
	"github.com/owezzy/schoolCRM/app/sdk/mid"
	"github.com/owezzy/schoolCRM/business/domain/admissionsbus"
	"github.com/owezzy/schoolCRM/foundation/web"
)

func authorizeAdmissions(authClient authclient.Authenticator, admissionsBus admissionsbus.ExtBusiness, rule string) web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, r *http.Request) web.Encoder {
			userID, err := mid.GetUserID(ctx)
			if err != nil {
				return errs.New(errs.Unauthenticated, err)
			}

			profile, err := admissionsBus.QueryStaffProfileByUserID(ctx, userID)
			if err != nil {
				if errors.Is(err, admissionsbus.ErrStaffProfileNotFound) {
					return errs.New(errs.Unauthenticated, admissionsbus.ErrStaffProfileNotFound)
				}
				return errs.Errorf(errs.Unauthenticated, "query admissions staff profile: %s", err)
			}

			if !profile.Active {
				return errs.New(errs.Unauthenticated, fmt.Errorf("admissions staff profile inactive"))
			}

			claims := mid.GetClaims(ctx)
			claims.Roles = append(claims.Roles, admissionsbus.AdmissionsRolesToStrings(profile.Roles)...)

			auth := authclient.Authorize{
				Claims: claims,
				UserID: userID,
				Rule:   rule,
			}

			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			if err := authClient.Authorize(ctx, auth); err != nil {
				return errs.New(errs.Unauthenticated, err)
			}

			return next(ctx, r)
		}

		return h
	}

	return m
}
