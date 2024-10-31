package property

import (
	"bytes"
	"fmt"
)

type Default struct{}

func (d *Default) Name() string {
    return "Default"
}

func (d *Default) Validate(diff PropertyDiff) (bool, error) {
	reset := func(diff PropertyDiff) PropertyDiff {
		oldProperty := diff.Old()
		newProperty := diff.New()
		oldProperty.Default = nil
		newProperty.Default = nil
		return NewPropertyDiff(oldProperty, newProperty)
	}

	var err error

	switch {
	case diff.Old().Default == nil && diff.New().Default != nil:
		err = fmt.Errorf("default value %q added when there was no default previously", string(diff.New().Default.Raw))
	case diff.Old().Default != nil && diff.New().Default == nil:
		err = fmt.Errorf("default value %q removed", string(diff.Old().Default.Raw))
	case diff.Old().Default != nil && diff.New().Default != nil && !bytes.Equal(diff.Old().Default.Raw, diff.New().Default.Raw):
		err = fmt.Errorf("default value changed from %q to %q", string(diff.Old().Default.Raw), string(diff.New().Default.Raw))
	}

	return IsHandled(diff, reset), err
}