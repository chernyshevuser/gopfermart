package business

func IsFinalized(status string) bool {
	return status == string(StatusInvalid) || status == string(StatusProcessed)
}

func GetInitialStatus() string {
	return string(StatusNew)
}

func GetFinalizedStatuses() []string {
	return []string{
		string(StatusInvalid),
		string(StatusProcessed),
	}
}

func GetNotFinalizedStatuses() []string {
	return []string{
		string(StatusNew),
		string(StatusProcessing),
		string(StatusRegistered),
	}
}

func GetAllStatuses() []string {
	finalized := GetFinalizedStatuses()
	notFinalized := GetNotFinalizedStatuses()

	res := make([]string, 0, len(finalized)+len(notFinalized))

	res = append(res, finalized...)
	res = append(res, notFinalized...)

	return res
}
