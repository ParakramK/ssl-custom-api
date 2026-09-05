package mysql

import (
	"context"
	"database/sql/driver"
	"errors"
	"testing"
)

type stubConnector struct {
	conn driver.Conn
	err  error
}

func (s stubConnector) Connect(context.Context) (driver.Conn, error) {
	return s.conn, s.err
}

func (s stubConnector) Driver() driver.Driver {
	return nil
}

type execStubConn struct {
	queries []string
	execErr error
	closed  bool
}

func (c *execStubConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("not implemented")
}

func (c *execStubConn) Close() error {
	c.closed = true
	return nil
}

func (c *execStubConn) Begin() (driver.Tx, error) {
	return nil, errors.New("not implemented")
}

func (c *execStubConn) ExecContext(_ context.Context, query string, _ []driver.NamedValue) (driver.Result, error) {
	c.queries = append(c.queries, query)
	return nil, c.execErr
}

type plainStubConn struct {
	closed bool
}

func (c *plainStubConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("not implemented")
}

func (c *plainStubConn) Close() error {
	c.closed = true
	return nil
}

func (c *plainStubConn) Begin() (driver.Tx, error) {
	return nil, errors.New("not implemented")
}

func TestReadonlyConnectorEnforcesReadOnly(t *testing.T) {
	inner := &execStubConn{}
	wrapped := readonlyConnector{Connector: stubConnector{conn: inner}}

	conn, err := wrapped.Connect(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if conn != driver.Conn(inner) {
		t.Error("expected the wrapped connection back")
	}

	if len(inner.queries) != 1 || inner.queries[0] != "SET SESSION TRANSACTION READ ONLY" {
		t.Errorf("expected read-only enforcement statement, got %v", inner.queries)
	}
}

func TestReadonlyConnectorPropagatesConnectError(t *testing.T) {
	wrapped := readonlyConnector{Connector: stubConnector{err: errors.New("dial failure")}}

	if _, err := wrapped.Connect(context.Background()); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestReadonlyConnectorRejectsNonExecConn(t *testing.T) {
	inner := &plainStubConn{}
	wrapped := readonlyConnector{Connector: stubConnector{conn: inner}}

	if _, err := wrapped.Connect(context.Background()); err == nil {
		t.Fatal("expected error, got nil")
	}

	if !inner.closed {
		t.Error("expected the unusable connection to be closed")
	}
}

func TestReadonlyConnectorClosesOnEnforcementFailure(t *testing.T) {
	inner := &execStubConn{execErr: errors.New("read-only unsupported")}
	wrapped := readonlyConnector{Connector: stubConnector{conn: inner}}

	if _, err := wrapped.Connect(context.Background()); err == nil {
		t.Fatal("expected error, got nil")
	}

	if !inner.closed {
		t.Error("expected the failed connection to be closed")
	}
}
