package main

func fmtBlack(s string) string {
	return "\033[30;1m" + s + "\033[0m"
}

func fmtYellow(s string) string {
	return "\033[33;1m" + s + "\033[0m"
}

func fmtGreen(s string) string {
	return "\033[32;1m" + s + "\033[0m"
}

func fmtRed(s string) string {
	return "\033[31;1m" + s + "\033[0m"
}
