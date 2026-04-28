package cli

const lowercase = "abcdefghijklmnopqrstuvwxyz"
const uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

const numbers = "1234567890"

const symbols = "-_"

func Alphabet() string {
	return lowercase + uppercase
}
