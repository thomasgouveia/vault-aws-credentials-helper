package structtoflags

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/spf13/cobra"
)

const (
	flagName        = "flag.name"
	flagDefault     = "flag.default"
	flagDescription = "flag.desc"
)

// MapStructToCommandFlags is a helper generic function that takes in parameter
// a struct and try to map the fields as command flags.
func MapStructToCommandFlags[T interface{}](cmd *cobra.Command, v *T) error {
	t := reflect.ValueOf(*v).Elem().Type()

	// Iterate over all fields in the given struct definition
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag
		var fName, fDesc, fDefault string

		fName, ok := tag.Lookup(flagName)
		if !ok {
			return errors.New("missing flag.name in struct field")
		}

		fDesc, ok = tag.Lookup(flagDescription)
		if !ok {
			fDesc = ""
		}

		fDefault, ok = tag.Lookup(flagDefault)
		if !ok {
			fDefault = ""
		}

		cmd.Flags().String(fName, fDefault, fDesc)
	}

	return nil
}

// MapStructToCommandFlags.
func MapCommandFlagsToStruct[T interface{}](cmd *cobra.Command, v *T) error {
	typeOf := reflect.ValueOf(*v).Elem().Type()
	valueOf := reflect.ValueOf(*v)

	// Iterate over all fields in the struct definition
	for i := 0; i < typeOf.NumField(); i++ {
		f := typeOf.Field(i)
		// Get the value of the struct to be able to
		// get the mutable field from it.
		obj := reflect.Indirect(valueOf)
		field, tag := obj.FieldByName(f.Name), f.Tag

		// Ensure the flag.name key is set into the field struct.
		// Otherwise, we should return an error
		fName, ok := tag.Lookup(flagName)
		if !ok {
			return fmt.Errorf("missing tag '%s' on field '%s' of struct %T for reflection", flagName, f.Name, v)
		}

		// Get the flag value
		val, err := cmd.Flags().GetString(fName)
		if err != nil {
			return err
		}

		// Set the value of the field
		field.SetString(val)
	}

	return nil
}
