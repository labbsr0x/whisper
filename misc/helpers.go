package misc

// ConvertInterfaceArrayToStringArray a helper method to perform conversions between []interface{} and []string
func ConvertInterfaceArrayToStringArray(toConvert []interface{}) []string {
	toReturn := make([]string, 0)
	if l := len(toConvert); l > 0 {
		toReturn = make([]string, l)
		for i := 0; i < l; i++ {
			toReturn[i] = toConvert[i].(string)
		}
	}
	return toReturn
}
