package validation

import "regexp"

// IsValidPhoneNumber checks if the phone number has a valid format.
// This regex matches phone numbers that are between 10 and 15 digits long.
func IsValidPhoneNumber(phone string) bool {
	// Regex for phone numbers with 10 to 15 digits
	const phoneRegex = `^\d{10,15}$`
	re := regexp.MustCompile(phoneRegex)
	return re.MatchString(phone)
}

// IsValidEmail checks if the provided email address is valid using a regular expression.
func IsValidEmail(email string) bool {
	// Simple regex to check the email format
	const emailRegex = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}
