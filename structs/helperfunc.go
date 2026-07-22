package extras

import (
	"encoding/gob"
	"unicode"
	"unicode/utf8"
)

func hasSpecialCharacter(s string) bool {
	for _, r := range s {
		if unicode.IsPunct(r) || unicode.IsSymbol(r) {
			return true
		}
	}
	return false
}
func hasUpperCase(s string) bool {
	for _, r := range s {
		if unicode.IsUpper(r) {
			return true
		}
	}
	return false
}
func hasLowerCase(s string) bool {
	for _, r := range s {
		if unicode.IsLower(r) {
			return true
		}
	}
	return false
}
func userNameAndPasswordAreValid(password string) bool {

	if utf8.RuneCountInString(password) < 16 {
		//error
		return false

	}
	// if !hasSpecialCharacter(password) {
	// 	return false

	// }
	if !hasUpperCase(password) {
		return false

	}
	if !hasLowerCase(password) {
		return false

	}
	return true

}

func RegisterEncodingDecodingTypes() {

	gob.Register(&SignUpAttempt{})
	gob.Register(&LoginAttempt{})
	gob.Register(&Message{})

}
