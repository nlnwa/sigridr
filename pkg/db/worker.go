package db

func SaveSearchResult(document interface{}) (string, error) {
	Connect()
	defer Disconnect()

	return Insert("results", document)
}
