package model

import (
	"fmt"
	"net/mail"
	"testing"
)

func TestParseModel(t *testing.T) {
	type args struct {
		id                string
		email             string
		phoneNumber       string
		firstName         string
		lastName          string
		consentConditions string
		consentMarketing  string
	}
	testcases := []struct {
		args    args
		wantErr bool
		want    Contestant
	}{
		{
			args: args{
				id:                "",
				email:             "example@example.com",
				phoneNumber:       "+48123456789",
				firstName:         "John",
				lastName:          "Doe",
				consentConditions: "on",
			},
			wantErr: false,
			want: Contestant{
				Email:             mail.Address{Address: "example@example.com"},
				PhoneNumber:       PhoneNumber{Value: "+48123456789"},
				FirstName:         "John",
				LastName:          "Doe",
				ConsentConditions: true,
				ConsentMarketing:  false,
			},
		},
	}

	for idx, tt := range testcases {
		t.Run(fmt.Sprintf("test_%d", idx), func(t *testing.T) {
			got, err := ParseContestant(
				tt.args.id,
				tt.args.email,
				tt.args.phoneNumber,
				tt.args.firstName,
				tt.args.lastName,
				tt.args.consentConditions,
				tt.args.consentMarketing,
			)
			if err != nil && !tt.wantErr {
				t.Errorf("got unexpected error %v", err)
			}
			if err == nil && tt.wantErr {
				t.Error("wanted error but go nil")
			}
			if !tt.wantErr && err == nil {
				if got.ID == "" ||
					got.Email.Address != tt.want.Email.Address ||
					got.PhoneNumber.Value != tt.want.PhoneNumber.Value ||
					got.FirstName != tt.want.FirstName ||
					got.LastName != tt.want.LastName ||
					got.ConsentConditions != tt.want.ConsentConditions ||
					got.ConsentMarketing != tt.want.ConsentMarketing {
					t.Errorf("got %v, want %v", got, tt.want)
				}
			}
			fmt.Println(got)
		})

	}
}
