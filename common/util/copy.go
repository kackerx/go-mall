package util

import (
	"errors"
	"regexp"
	"time"

	"github.com/jinzhu/copier"
)

func Copy(dst, src any) error {
	return copier.CopyWithOption(dst, src, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
		Converters: []copier.TypeConverter{
			{
				SrcType: time.Time{},
				DstType: copier.String,
				Fn: func(src interface{}) (dst interface{}, err error) {
					s, ok := src.(time.Time)
					if !ok {
						return nil, errors.New("src type not time")
					}

					return s.Format(time.DateTime), nil
				},
			},
			{
				SrcType: copier.String,
				DstType: time.Time{},
				Fn: func(src interface{}) (dst interface{}, err error) {
					s, ok := src.(string)
					if !ok {
						return nil, errors.New("src type not string")
					}

					pattern := `^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$`
					matched, _ := regexp.MatchString(pattern, s)
					if !matched {
						return nil, errors.New("src not format string")
					}

					return time.Parse(time.DateTime, s)
				},
			},
		},
	})
}
