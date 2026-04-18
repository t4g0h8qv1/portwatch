package portprofile

// Default returns a Registry pre-loaded with common host-role profiles.
func Default() *Registry {
	r := New()
	_ = r.Register("web", []int{80, 443, 8080, 8443})
	_ = r.Register("database", []int{3306, 5432, 1433, 27017, 6379})
	_ = r.Register("ssh", []int{22})
	_ = r.Register("mail", []int{25, 465, 587, 110, 995, 143, 993})
	_ = r.Register("dns", []int{53})
	_ = r.Register("monitoring", []int{9090, 9100, 3000, 8086})
	return r
}
