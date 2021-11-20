package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/jackc/pgx/v4"

	pb "github.com/Serj1c/usermngnt/proto"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

// UserManagementServer ...
type UserManagementServer struct {
	db *pgx.Conn
	pb.UnimplementedUserManagementServer
}

// NewUserManagementServer is a constructor for UserManagementServer
func NewUserManagementServer(db *pgx.Conn) *UserManagementServer {
	return &UserManagementServer{
		db: db,
	}
}

// CreateNewUser creates a new user and appends it the users list
func (ums *UserManagementServer) CreateNewUser(ctx context.Context, in *pb.NewUser) (*pb.User, error) {

	sql := `CREATE TABLE IF NOT EXISTS users(
		id SERIAL PRIMARY KEY,
		name VARCHAR,
		age int
	);`
	_, err := ums.db.Exec(context.Background(), sql)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Table creation failed: %v\n", err)
		os.Exit(1)
	}
	createdUser := &pb.User{
		Name: in.GetName(),
		Age:  in.GetAge(),
	}
	tx, err := ums.db.Begin(context.Background())
	defer tx.Rollback(context.Background())
	if err != nil {
		fmt.Printf("Start of transaction failed: %s", err)
	}
	_, err = tx.Exec(context.Background(), "INSERT INTO users(name, age) values($1, $2)", createdUser.Name, createdUser.Age)
	if err != nil {
		fmt.Printf("tx.Exec failed: %v", err)
	}
	tx.Commit(context.Background())

	return createdUser, nil
}

// GetUsers returns a list of users
func (ums *UserManagementServer) GetUsers(ctx context.Context, in *pb.GetUsersParams) (*pb.UsersList, error) {
	userList := &pb.UsersList{}
	rows, err := ums.db.Query(context.Background(), "SELECT id, name, age FROM users")
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		user := pb.User{}
		err = rows.Scan(&user.Id, &user.Name, &user.Age)
		if err != nil {
			return nil, err
		}
		userList.Users = append(userList.Users, &user)
	}
	return userList, nil
}

// Run starts the server
func (ums *UserManagementServer) Run() error {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	pb.RegisterUserManagementServer(s, ums)
	log.Printf("server is listening on %v", listener.Addr())
	return s.Serve(listener)
}

func main() {

	dbURL := "postgres://postgres:password@localhost:5432/postgres"
	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatal("unable to get a db connection")
	}
	defer conn.Close(context.Background())

	userManSer := NewUserManagementServer(conn)
	if err := userManSer.Run(); err != nil {
		log.Fatal(err)
	}
}
