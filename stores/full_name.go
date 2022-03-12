package stores

func GenerateFullName() string {
	return GetRandomFirstName() + GetRandomMiddleName()
}
