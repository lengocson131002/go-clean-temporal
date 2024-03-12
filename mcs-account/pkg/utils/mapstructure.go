package utils

import "github.com/mitchellh/mapstructure"

func MapStruct(input interface{}, output interface{}) error {
	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           &output,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	err = decoder.Decode(input)
	return err
}
