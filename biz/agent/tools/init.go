package tools

func InitTools() {
	err := initSlotTools()
	if err != nil {
		panic(err)
	}
	err = InitCommonTools()
	if err != nil {
		panic(err)
	}
}
