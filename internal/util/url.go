package util

import "net/url"

func URLValuesWith(values url.Values, key string, value string) url.Values {
	newValues := url.Values{}
	for key, valueSet := range values {
		for _, value := range valueSet {
			newValues.Add(key, value)
		}
	}

	newValues.Set(key, value)
	return newValues
}

func URLValuesWithout(values url.Values, keyWithout string) url.Values {
	newValues := url.Values{}
	for key, valueSet := range values {
		if key == keyWithout {
			continue
		}
		for _, value := range valueSet {
			newValues.Add(key, value)
		}
	}

	return newValues
}
