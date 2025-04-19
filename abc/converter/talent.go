package converter

import (
	domain "g123-jp/talent/app/domain/talent"
	infra "g123-jp/talent/app/infrastructure/talent"
)

// goverter:output:file @cwd/app/infrastructure/talent/convert.go
// goverter:converter
// goverter:output:format function
// goverter:extend g123-jp/talent/pkg/util/conv:Convert.*
// goverter:ignoreMissing
type TalentForInfra interface {
	// goverter:useZeroValueOnPointerInconsistency
	// goverter:ignoreMissing
	ToModel(do *domain.Talent) *infra.Talent
}
