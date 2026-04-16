package portgroup

// WellKnown returns a Registry pre-loaded with common service port groups.
func WellKnown() *Registry {
	r := New()
	defs := []struct {
		name  string
		ports []int
	}{
		{"web", []int{80, 443, 8080, 8443}},
		{"ssh", []int{22}},
		{"dns", []int{53}},
		{"mail", []int{25, 465, 587, 993, 995}},
		{"database", []int{3306, 5432, 1433, 27017, 6379}},
		{"ftp", []int{20, 21}},
		{"ldap", []int{389, 636}},
		{"monitoring", []int{9090, 9100, 9200, 9300}},
	}
	for _, d := range defs {
		// Errors are impossible with hard-coded valid values.
		_ = r.Register(d.name, d.ports)
	}
	return r
}
