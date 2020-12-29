package collection

// Diff returns the additions and removals to/from r, compared to l
func Diff(l []interface{}, r []interface{}) (additions []interface{}, removals []interface{}) {
	additions = Additions(l, r)
	// Reversing the direction gives us the removals
	removals = Additions(r, l)

	return
}

// Additions returns a list of items that exist in n, but not in o
func Additions(l []interface{}, r []interface{}) (additions []interface{}) {
	// Check which values have been additions
	for _, rv := range r {
		if !Contains(rv, l) {
			additions = append(additions, rv.(string))
		}
	}

	return
}

// Contains returns true when the haystack contains the needle
func Contains(needle interface{}, haystack []interface{}) bool {
	for _, hay := range haystack {
		if needle == hay {
			return true
		}
	}
	return false
}

// Ordered returns all the items from r, while retaining the ordering of l for existing items
func Ordered(items []interface{}, like []interface{}) (ordered []interface{}) {
	// First, add items that exist in like, retaining their order
	for _, l := range like {
		if Contains(l, items) {
			ordered = append(ordered, l)
		}
	}

	// Next, add items that don't exist in like next,using the order in which they appear
	for _, item := range items {
		if !Contains(item, like) {
			ordered = append(ordered, item)
		}
	}

	return
}

// ToInterfaceSlice returns the string slice as an interface slice
func ToInterfaceSlice(s []string) []interface{} {
	output := make([]interface{}, len(s))
	for i, j := range s {
		output[i] = j
	}
	return output
}
