package models_test

import (
	"testing"

	"github.com/DiegoTUI/signpost/models"
)

func TestArrayExtract(t *testing.T) {

	array := []*models.Sign{
		&models.Sign{
			Distance: 0,
		},
		&models.Sign{
			Distance: 1,
		},
		&models.Sign{
			Distance: 2,
		},
		&models.Sign{
			Distance: 3,
			Angle:    2,
		},
	}

	if _, _, err := models.SignArrayExtract(array, 4); err == nil {
		t.Error("ArrayExtract did not return index out of bounds")
	}

	// extract element in the middle
	item, sliced, err := models.SignArrayExtract(array, 2)
	if err != nil {
		t.Error("ArrayExtract did return an error for element 2")
	}
	if item.Distance != 2 {
		t.Error("ArrayExtract returned wrong item for element 2")
	}
	if len(sliced) != 3 ||
		sliced[0].Distance != 0 ||
		sliced[1].Distance != 1 ||
		sliced[2].Distance != 3 {
		t.Error("ArrayExtract returned wrong slice for element 2")
	}
	if len(array) != 4 ||
		array[0].Distance != 0 ||
		array[1].Distance != 1 ||
		array[2].Distance != 2 ||
		array[3].Distance != 3 {
		t.Error("ArrayExtract modified original array for element 2")
	}

	// extract last element
	item, sliced, err = models.SignArrayExtract(array, 3)
	if err != nil {
		t.Error("ArrayExtract did return an error for element 3")
	}
	if item.Distance != 3 {
		t.Error("ArrayExtract returned wrong item for element 3")
	}
	if len(sliced) != 3 ||
		sliced[0].Distance != 0 ||
		sliced[1].Distance != 1 ||
		sliced[2].Distance != 2 {
		t.Error("ArrayExtract returned wrong slice for element 3")
	}
	if len(array) != 4 ||
		array[0].Distance != 0 ||
		array[1].Distance != 1 ||
		array[2].Distance != 2 ||
		array[3].Distance != 3 {
		t.Error("ArrayExtract modified original array for element 3")
	}

	// extract first element
	item, sliced, err = models.SignArrayExtract(array, 0)
	if err != nil {
		t.Error("ArrayExtract did return an error for element 0")
	}
	if item.Distance != 0 {
		t.Error("ArrayExtract returned wrong item for element 0")
	}
	if len(sliced) != 3 ||
		sliced[0].Distance != 1 ||
		sliced[1].Distance != 2 ||
		sliced[2].Distance != 3 {
		t.Error("ArrayExtract returned wrong slice for element 0")
	}
	if len(array) != 4 ||
		array[0].Distance != 0 ||
		array[1].Distance != 1 ||
		array[2].Distance != 2 ||
		array[3].Distance != 3 {
		t.Error("ArrayExtract modified original array for element 0")
	}
}
