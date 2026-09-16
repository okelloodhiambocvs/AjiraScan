package auth

import "fmt"

type Role string

const (
	PlatformAdmin     Role = "PLATFORM_ADMIN"
	OrganizationAdmin Role = "ORGANIZATION_ADMIN"
	HRRecruiter       Role = "HR_RECRUITER"
	HiringManager     Role = "HIRING_MANAGER"
	Applicant         Role = "APPLICANT"
)

func (r Role) Valid() bool {
	switch r {
	case PlatformAdmin, OrganizationAdmin, HRRecruiter, HiringManager, Applicant:
		return true
	default:
		return false
	}
}

func Authorize(actual Role, permitted ...Role) error {
	if actual == PlatformAdmin {
		return nil
	}
	for _, role := range permitted {
		if actual == role {
			return nil
		}
	}
	return fmt.Errorf("forbidden for role %s", actual)
}
