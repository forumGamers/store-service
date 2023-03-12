package helper

func QueryBuild(init string, add string) string {
	query := ""
	if init == "" {
		query = add
	} else {
		query = init + " and " + add
	}
	return query
}