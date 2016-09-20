package main

import "testing"

func Test_if_not_able_to_parse_return_nil(t *testing.T) {
	_, err := Parse("kksajfs")
	if err == nil {
		t.Errorf("Expected an error value: %v", err)
	}
}

func Test_single_name_is_parsed_from_xml(t *testing.T) {
	data := `<h2><a href="/grafton-east-car-park">Grafton East car park</a></h2>
<p><strong>648 spaces</strong> (17% full and emptying)</p>`
	expName := "Grafton East car park"
	expSpaces := 648
	d, err := Parse(data)
	if err != nil {
		t.Errorf("Expected an error value: %v", err)
	}
	if len(d) != 1 {
		t.Error("Expected exactly one element in slice")
	}
	if d[0].Name != expName {
		t.Errorf("Expected %s, got '%s'", expName, d[0].Name)
	}
	if d[0].Spaces != expSpaces {
		t.Errorf("Expected %d, got '%d'", expSpaces, d[0].Spaces)
	}
}

func Test_a_complete_set_of_data_can_be_parsed(t *testing.T) {
	data := `<h2><a href="/grafton-east-car-park">Grafton East car park</a></h2>
	<p><strong>648 spaces</strong> (17% full and emptying)</p>
	<h2><a href="/grafton-west-car-park">Grafton West car park</a></h2>
	<p><strong>132 spaces</strong> (53% full and filling)</p>
	<h2><a href="/grand-arcade-car-park">Grand Arcade car park</a></h2>
	<p><strong>190 spaces</strong> (79% full and filling)</p>
	<h2><a href="/park-street-car-park">Park Street car park</a></h2>
	<p><strong>181 spaces</strong> (52% full and filling)</p>
	<h2><a href="/queen-anne-terrace-car-park">Queen Anne Terrace car park</a></h2>
	<p><strong>285 spaces</strong> (47% full and filling)</p>`
	d, err := Parse(data)
	if err != nil {
		t.Errorf("Expected an error value: %v", err)
	}
	if len(d) != 5 {
		t.Errorf("Expected exactly %d element in slice, got %d", 5, len(d))
	}
	expected := []struct {
		n string
		v int
	}{
		{"Grafton East car park", 648},
		{"Grafton West car park", 132},
		{"Grand Arcade car park", 190},
		{"Park Street car park", 181},
		{"Queen Anne Terrace car park", 285},
	}
	for i, e := range expected {
		if d[i].Name != e.n {
			t.Errorf("Name error, got %s, expected %s", d[i].Name, e.n)
		}
		if d[i].Spaces != e.v {
			t.Errorf("Name error, got %d, expected %d", d[i].Spaces, e.v)
		}
	}
}
