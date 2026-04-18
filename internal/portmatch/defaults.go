package portmatch

// Default returns a Matcher pre-loaded with common service rule sets.
func Default() *Matcher {
	m := New()
	// errors ignored: all values are valid constants
	_ = m.Add("web", "HTTP/HTTPS services", []int{80, 443, 8080, 8443})
	_ = m.Add("ssh", "Secure Shell", []int{22})
	_ = m.Add("database", "Common database ports", []int{3306, 5432, 1433, 27017, 6379})
	_ = m.Add("mail", "Mail transfer and retrieval", []int{25, 465, 587, 110, 143, 993, 995})
	_ = m.Add("dns", "Domain Name System", []int{53})
	_ = m.Add("ftp", "File Transfer Protocol", []int{20, 21})
	_ = m.Add("monitoring", "Metrics and monitoring", []int{9090, 9100, 3000, 8086})
	return m
}
