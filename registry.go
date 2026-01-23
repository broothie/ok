package ok

func Registry() []Tool {
	return []Tool{
		NewMake(),
		NewNPM(),
	}
}
