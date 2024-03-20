package validators

import (
	"regexp"
	"strings"
)

func ValidateMobilePhone(phone string) (string, bool) {
	if !IsValidPhone(phone) {
		return phone, false
	}
	trimmedPhone := trimMobilePhone(phone)

	//if not is kazakhstan and russian
	if string(trimmedPhone[0]) != "7" && string(trimmedPhone[0]) != "8" {
		return "+" + trimmedPhone, true
	}

	if len(trimmedPhone) != 10 && len(trimmedPhone) != 11 {
		return phone, false
	}
	//russian number
	if trimmedPhone[:2] == "89" || trimmedPhone[:2] == "79" {
		return "+" + trimmedPhone, true
	}

	//validation KazakhstanPhone
	if !IsKazakhstanPhone(trimmedPhone) {
		return phone, true
	}

	if IsKazakhstanPhone(trimmedPhone) {
		return "+" + transformFirstSymbol(trimmedPhone), true
	}

	return "+" + trimmedPhone, true
}

// IsValidPhone Проверяет чтоб были допустимые символы и первый знак +
func IsValidPhone(str string) bool {

	re := regexp.MustCompile(`^[+]?[0-9()\s–-]+$`)
	return re.MatchString(str)
}

func isFirstPlus(phone string) bool {
	if phone[0] == '+' {
		return true
	}
	return false
}

func transformFirstSymbol(phone string) string {
	if len(phone) == 11 {
		if []rune(phone)[0] == '7' || []rune(phone)[0] == '8' {
			newPhone := []rune(phone)
			newPhone[0] = '7'
			return string(newPhone)
		}
	}
	if len(phone) == 10 {
		return "7" + string(phone)
	}

	return phone
}
func trimMobilePhone(input string) string {
	// Символы, которые нужно удалить
	charsToRemove := "()-–-+b "

	// Итерируем по строке и удаляем символы
	for _, char := range charsToRemove {
		input = strings.ReplaceAll(input, string(char), "")
	}

	return input
}

func IsKazakhstanPhone(phone string) bool {
	operatorCodes := []string{"701", "702", "703", "704", "705",
		"706", "707", "708", "709", "747", "750", "751", "760", "761",
		"762", "763", "764", "771", "775", "776", "777", "778", "301",
		"302", "303", "304", "305", "306", "307", "308", "309", "347",
		"350", "351", "360", "361", "362", "363", "364", "371", "375",
		"376", "377", "378"}

	if len(phone) == 11 && (string(phone[0]) == "8" || string(phone[0]) == "7") {
		code := phone[1:4]
		exist := false
		for _, s := range operatorCodes {
			if s == code {
				exist = true
			}
		}
		if exist {
			return true
		}
	}

	if len(phone) == 10 {
		code := phone[0:3]
		exist := false
		for _, s := range operatorCodes {
			if s == code {
				exist = true
			}
		}
		if exist {
			return true
		}
	}

	return false
}

func ValidateBasicPhone(phone string) (string, bool) {
	if !IsValidPhone(phone) {
		return phone, false
	}
	firstPlus := isFirstPlus(phone)
	trimmedPhone := trimMobilePhone(phone)

	if firstPlus {
		trimmedPhone = "+" + trimmedPhone
	}
	return trimmedPhone, true
}
