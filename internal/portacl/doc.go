// Package portacl provides a simple ordered access-control list (ACL)
// for port numbers.
//
// Rules are evaluated in insertion order; the first matching rule
// determines the Action (Allow or Deny). If no rule matches, the
// default action is Allow.
//
// Example:
//
//	acl := portacl.New()
//	acl.Add(portacl.Rule{Port: 22, Direction: portacl.Inbound, Action: portacl.Deny})
//	action := acl.Evaluate(22, portacl.Inbound) // Deny
package portacl
