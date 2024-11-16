package object

func NativeToInteger(native int) *Integer {
	return &Integer{Value: native}
}

func NativeToBoolean(native bool) *Boolean {
	if native {
		return TRUE
	}
	return FALSE
}

func NativeToString(native string) *String {
	return &String{Value: native}
}
