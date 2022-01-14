package main

import (
	"testing"
)

func Test_storeTxnInfo(t *testing.T) {
	type args struct {
		txn    *TransactionInfo
		txn_id int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test case: 1",
			args: args{
				txn: &TransactionInfo{
					Amount:    100.0,
					Txn_type:  "pen",
					Parent_id: 0,
				},
				txn_id: 1,
			},
			wantErr: false,
		},
		{
			name: "test case: 2",
			args: args{
				txn: &TransactionInfo{
					Amount:    200.50,
					Txn_type:  "pen",
					Parent_id: 2,
				},
				txn_id: 3,
			},
			wantErr: true,
		},
		{
			name: "test case: 3",
			args: args{
				txn: &TransactionInfo{
					Amount:    300.0,
					Txn_type:  "car",
					Parent_id: 1,
				},
				txn_id: 5,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := storeTxnInfo(tt.args.txn, tt.args.txn_id); (err != nil) != tt.wantErr {
				t.Errorf("storeTxnInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_calSum(t *testing.T) {
	type args struct {
		txn_id int64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		// TODO: Add test cases.
		{
			name: "sum_test: 1",
			args: args{
				txn_id: 1,
			},
			want: 400.0,
		},
		{
			name: "sum_test: 2",
			args: args{
				txn_id: 5,
			},
			want: 300.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calSum(tt.args.txn_id); got != tt.want {
				t.Errorf("calSum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_deleteTxn(t *testing.T) {
	type args struct {
		txn_id int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "delete_test: 1",
			args: args{
				txn_id: 1,
			},
			wantErr: true,
		},
		{
			name: "delete_test: 2",
			args: args{
				txn_id: 4,
			},
			wantErr: true,
		},
		{
			name: "delete_test: 3",
			args: args{
				txn_id: 5,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := deleteTxn(tt.args.txn_id); (err != nil) != tt.wantErr {
				t.Errorf("deleteTxn() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
