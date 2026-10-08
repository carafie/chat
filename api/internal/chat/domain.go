package chat

import (
	"net/url"
	"unicode"

	"github.com/carafie/chat/api/pkg/errcode"
	"github.com/carafie/chat/api/pkg/text"
)

var (
	ErrDomainNameTooShort = errcode.New("DOMAIN_NAME_TOO_SHORT", "domain name is too short")
	ErrDomainNameTooLong  = errcode.New("DOMAIN_NAME_TOO_LONG", "domain name is too long")
	ErrDomainNameInvalid  = errcode.New("DOMAIN_NAME_INVALID", "domain name is invalid")
)

type Domain struct {
	Name string
}

func NewDomain(name string) (Domain, error) {
	url, err := url.Parse(name)
	if err != nil {
		return Domain{}, ErrDomainNameInvalid
	}
	name, err = text.Parser{
		MinChars: 1,
		MaxChars: 255,
		CharValid: func(char rune) bool {
			return unicode.IsPrint(char) && !unicode.IsSpace(char)
		},
	}.Parse(url.Hostname())
	switch err {
	case nil:
		return Domain{Name: name}, nil
	case text.ErrTooShort:
		return Domain{}, ErrDomainNameTooShort
	case text.ErrTooLong:
		return Domain{}, ErrDomainNameTooLong
	default:
		return Domain{}, ErrDomainNameInvalid
	}
}
