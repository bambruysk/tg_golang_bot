package main

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4"
)

func TestConnectDB(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"simple"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ConnectDB()
		})
	}
}

func TestUser_ReadFromDB(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()
	type fields struct {
		State      DialogState
		CurrHolde  int
		CurrPlayer *Player
		Name       string
		Location   string
	}
	type args struct {
		chat_id int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "no user",
			fields: fields{
				State:      IDLE,
				CurrHolde:  -1,
				CurrPlayer: nil,
				Name:       "Test user",
				Location:   "Test location",
			},
			args: args{
				chat_id: 13,
			},
			wantErr: false,
		},
	}
	conn, err := ConnectDB()
	if err != nil {
		t.Errorf("Connet to tatabse failed %v", err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				State:      tt.fields.State,
				CurrHolde:  tt.fields.CurrHolde,
				CurrPlayer: tt.fields.CurrPlayer,
				Name:       tt.fields.Name,
				Location:   tt.fields.Location,
			}
			if err := u.ReadFromDB(conn); (err != nil) != tt.wantErr {
				t.Errorf("User.ReadFromDB() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	t.Fail()
	t.Log(buf.String())
}

func TestLocation_AddDB(t *testing.T) {
	type fields struct {
		Name string
	}
	type args struct {
		conn *pgx.Conn
	}
	conn, err := ConnectDB()
	if err != nil {
		t.Fatal("Connect to database fail")
		return
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "init",
			fields: fields{
				Name: "Aktlan",
			},
			args: args{
				conn: conn,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Location{
				Name: tt.fields.Name,
			}
			if err := l.AddDB(tt.args.conn); (err != nil) != tt.wantErr {
				t.Errorf("Location.AddDB() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
