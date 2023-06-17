package helper

func QueryBuild(query *string, add string) {
	if *query == "" {
		*query = add
	} else {
		*query = *query + " and " + add
	}
}