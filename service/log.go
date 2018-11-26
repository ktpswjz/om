package main

func LogError(v ...interface{}) string {
	if log == nil {
		return ""
	}

	return log.Error(v...)
}

func LogInfo(v ...interface{}) string {
	if log == nil {
		return ""
	}

	return log.Info(v...)
}
