package main

import (
	"flag"
)

type Flags struct {
	email                string
	domain               string
	isBaseConfigComplete bool

	forceResetTemplates bool
}

var flags Flags

const (
	DomainDescription = "Specify Base Domain Name"
	EmailDescription  = "Specify Your Email for SSL certificate generation"
)

func init() {
	flag.BoolVar(&flags.forceResetTemplates, "f", false, "Reset default templates")
	flag.BoolVar(&flags.forceResetTemplates, "forceDefault", false, "Reset default templates")

	flag.StringVar(&flags.domain, "d", "", DomainDescription)
	flag.StringVar(&flags.domain, "domain", "", DomainDescription)

	flag.StringVar(&flags.email, "email", "", EmailDescription)

	flag.Parse()

	flags.isBaseConfigComplete = false
	if flags.domain != "" && flags.email != "" {
		flags.isBaseConfigComplete = true
	}
}

func getBaseConfigVariables() []*Variable {
	focusTaken := false
	var vars []*Variable

	if flags.email == "" {
		ti := buildTextInput("your@email.com")
		if !focusTaken {
			ti.Focus()
			focusTaken = true
		}

		vars = append(vars, &Variable{
			Input:       ti,
			Name:        "Email",
			Template:    "base-config",
			Default:     "your@email.com",
			Description: EmailDescription,
		})
	}

	if flags.domain == "" {
		ti := buildTextInput("localhost")
		if !focusTaken {
			ti.Focus()
			focusTaken = true
		}

		vars = append(vars, &Variable{
			Input:       ti,
			Name:        "Domain",
			Default:     "localhost",
			Template:    "base-config",
			Description: DomainDescription,
		})
	}

	return vars
}
