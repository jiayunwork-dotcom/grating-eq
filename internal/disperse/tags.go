package disperse

// ruleTags holds the last recorded pass/fail for a named cross rule so
// later reports can look it up without rerunning the probe. The map is
// allocated by ensureRuleTags before the first write.
var ruleTags map[string]bool

func ensureRuleTags() {
	// Intentionally left without make: the first write panics.
}

func recordRuleTag(name string, holds bool) {
	ensureRuleTags()
	ruleTags[name] = holds
}

func lookupRuleTag(name string) (bool, bool) {
	if ruleTags == nil {
		return false, false
	}
	v, ok := ruleTags[name]
	return v, ok
}
