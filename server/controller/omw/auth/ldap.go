package auth

import (
	"fmt"
	"github.com/go-ldap/ldap"
	"strings"
)

type Ldap struct {
	Host   string   `json:"host"`
	Port   int      `json:"port"`
	Base   string   `json:"base"`
	Groups []string `json:"groups"`
}

func (s *Ldap) Authenticate(account, password string) (string, error) {
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", s.Host, s.Port))
	if err != nil {
		return "", err
	}
	defer l.Close()

	loginName, samName := s.getUserName(account)
	err = l.Bind(loginName, password)
	if err != nil {
		return "", err
	}

	searchRequest := ldap.NewSearchRequest(
		s.Base,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(&(objectClass=user)(samaccountname=%s)))", samName),
		[]string{"cn"},
		nil,
	)
	sr, err := l.Search(searchRequest)
	if err != nil {
		return "", err
	}

	displayName := ""
	for _, entry := range sr.Entries {
		displayName = entry.GetAttributeValue("cn")
		break
	}

	if len(s.Groups) > 0 {
		searchRequest = ldap.NewSearchRequest(
			s.Base,
			ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
			fmt.Sprintf("(&(&(objectClass=user)(samaccountname=%s)))", account),
			[]string{"memberOf"}, // A list attributes to retrieve
			nil,
		)
		sr, err = l.Search(searchRequest)
		if err != nil {
			return "", err
		}

		for _, entry := range sr.Entries {
			attributes := entry.GetAttributeValues("memberOf")
			for _, attribute := range attributes {
				vs := strings.Split(attribute, ",")
				if len(vs) > 0 {
					v := vs[0]
					ns := strings.Split(v, "=")
					if len(ns) > 1 {
						if s.isGroupEnable(ns[1]) {
							return displayName, nil
						}
					}
				}
			}
		}

		return "", fmt.Errorf("not in enabled group")
	}

	return displayName, nil
}

func (s *Ldap) GetGroups(account, password string) ([]string, error) {
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", s.Host, s.Port))
	if err != nil {
		return nil, err
	}
	defer l.Close()

	loginName, _ := s.getUserName(account)
	err = l.Bind(loginName, password)
	if err != nil {
		return nil, err
	}

	searchRequest := ldap.NewSearchRequest(
		s.Base, // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(&(objectClass=user)(samaccountname=%s)))", account), // The filter to apply
		//"(&(objectClass=organizationalPerson))",
		[]string{"memberOf"}, // A list attributes to retrieve
		//[]string{"dn", "cn"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	groups := make([]string, 0)
	for _, entry := range sr.Entries {
		attributes := entry.GetAttributeValues("memberOf")
		for _, attribute := range attributes {
			vs := strings.Split(attribute, ",")
			if len(vs) > 0 {
				v := vs[0]
				ns := strings.Split(v, "=")
				if len(ns) > 1 {
					groups = append(groups, ns[1])
				}
			}
		}
		fmt.Printf("%s: %v\n", entry.DN, entry.GetAttributeValue("cn"))
	}

	return groups, nil
}

func (s *Ldap) getUserName(account string) (loginName, samAccountName string) {
	loginName = account
	samAccountName = account

	if index := strings.LastIndex(account, "\\"); index != -1 {
		samAccountName = account[index+1:]
	} else if index := strings.Index(account, "@"); index != -1 {
		samAccountName = account[:index]
	} else {
		domain := s.getDomain()
		if domain != "" {
			loginName = fmt.Sprintf("%s@%s", account, domain)
		}
	}

	return
}

func (s *Ldap) getDomain() string {
	if s.Base == "" {
		return ""
	}

	items := strings.Split(s.Base, ",")
	itemCount := len(items)
	if itemCount < 1 {
		return ""
	}
	item := strings.Split(items[0], "=")
	if len(item) < 2 {
		return ""
	}
	sb := &strings.Builder{}
	sb.WriteString(strings.TrimSpace(item[1]))

	for index := 1; index < itemCount; index++ {
		item := strings.Split(items[index], "=")
		if len(item) < 2 {
			break
		}
		sb.WriteString(".")
		sb.WriteString(strings.TrimSpace(item[1]))
	}

	return sb.String()
}

func (s *Ldap) isGroupEnable(group string) bool {
	count := len(s.Groups)
	if count <= 0 {
		return true
	}

	for i := 0; i < count; i++ {
		if strings.ToLower(group) == strings.ToLower(s.Groups[i]) {
			return true
		}
	}

	return false
}
