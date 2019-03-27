package obj

import "testing"

func TestAddressStr(t *testing.T) {
	a := &Address{
		Street: "1600 Pennsylvania Ave NW", CityName: "Washington",
		State: "dc", Zipcode: "20500",
	}
	expected := []string{
		`1600 Pennsylvania Ave NW
Washington, DC 20500`,
		`1600 Pennsylvania Ave NW
   Washington, DC 20500`,
	}
	formatted := []string{AddressFmt(a), AddressFmtIndent(a, 3)}
	for i, exp := range expected {
		if exp != formatted[i] {
			t.Errorf("unexpected output...\ngot:\n%s\nwanted:\n%s\n", formatted[i], exp)
		}
	}

	a = &Address{
		Street: "1600 Pennsylvania Ave NW", CityName: "Washington DC",
		State: "", Zipcode: "20500",
	}
	expected = []string{
		`1600 Pennsylvania Ave NW
Washington DC, 20500`,
		`1600 Pennsylvania Ave NW
   Washington DC, 20500`,
	}
	formatted = []string{AddressFmt(a), AddressFmtIndent(a, 3)}
	for i, exp := range expected {
		if exp != formatted[i] {
			t.Errorf("unexpected output...\ngot:\n%s\nwanted:\n%s\n", formatted[i], exp)
		}
	}
	t.Run("bad_zip_spaces", func(t *testing.T) {
		a.Zipcode = "        "
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected a panic")
			}
		}()
		_ = a.String()
	})
	t.Run("bad_zip_spaces", func(t *testing.T) {
		a.Zipcode = "123456"
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected a panic")
			}
		}()
		_ = a.String()
	})
	t.Run("bad_zip_spaces", func(t *testing.T) {
		a.Zipcode = "12345"
		a.State = "ABC"
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected a panic")
			}
		}()
		_ = a.String()
	})
}