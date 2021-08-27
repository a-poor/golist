package main

func toBlack(s string) string {
	return "\033[30;1m" + s + "\033[0m"
}

func toYellow(s string) string {
	return "\033[33;1m" + s + "\033[0m"
}

func toGreen(s string) string {
	return "\033[32;1m" + s + "\033[0m"
}

func toRed(s string) string {
	return "\033[31;1m" + s + "\033[0m"
}
