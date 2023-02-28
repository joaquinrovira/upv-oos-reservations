package daytime

import "testing"

func TestDayTime_Validate(t *testing.T) {
	type fields struct {
		Hour   uint8
		Minute uint8
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{name: "Default sctruct is valid", fields: fields{}, wantErr: false},
		{name: "Valid non-zero time", fields: fields{Hour: 12, Minute: 12}, wantErr: false},
		{name: "Invalid when hour > 23", fields: fields{Hour: 24}, wantErr: true},
		{name: "Invalid when minute > 59", fields: fields{Minute: 60}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := DayTime{
				Hour:   tt.fields.Hour,
				Minute: tt.fields.Minute,
			}
			if err := tr.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("DayTime.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
