package external

import (
	"fmt"
	)

type MockRedisConn struct {
	command_list []string 
}
func (*MockRedisConn) Close() error {
	return nil
}
func (*MockRedisConn) Err() error {
	return nil
}
func (conn *MockRedisConn) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	command := commandName
	for _, arg := range args {
		switch arg := arg.(type) {
		case string:
			command = command + " " + arg
		}
	}
	conn.command_list = append(conn.command_list, command)
	return nil, nil
}
func (conn *MockRedisConn) Send(commandName string, args ...interface{}) error {
	command := commandName
	for _, arg := range args {
		switch arg := arg.(type) {
		case string:
			command = command + " " + arg
		case uint32:
			command = fmt.Sprintf("%s %d", command, arg)
		}
	}
	conn.command_list = append(conn.command_list, command)
	return nil
}
func (*MockRedisConn) Flush() error {
	return nil
}
func (*MockRedisConn) Receive() (reply interface{}, err error) {
	return nil, nil
}