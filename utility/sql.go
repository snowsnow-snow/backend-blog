package utility

func OrByParams(field string, length int) string {
	var where string
	for i := 0; i < length; i++ {
		where += " " + field + " = ? "
		if (i + 1) < length {
			where += " or "
		}
	}
	return where
}
