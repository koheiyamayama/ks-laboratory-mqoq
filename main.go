package main

import (
	"context"
	"crypto/tls"
	"net"
	"os"
	"sync"

	"github.com/quic-go/quic-go"
	"golang.org/x/exp/slog"
)

func main() {
	ctx := context.Background()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))
	udpConn, err := net.ListenUDP("udp4", &net.UDPAddr{Port: 1234})
	if err != nil {
		logger.ErrorContext(ctx, err.Error())
		os.Exit(1)
	}
	tr := quic.Transport{
		Conn: udpConn,
	}
	ln, err := tr.Listen(&tls.Config{}, &quic.Config{})
	if err != nil {
		logger.ErrorContext(ctx, err.Error())
		os.Exit(1)
	}

	logger.InfoContext(ctx, "start ks-laboratory-mqoq on 1234 port")
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			conn, err := ln.Accept(ctx)
			if err != nil {
				logger.ErrorContext(ctx, err.Error())
				os.Exit(1)
			}
			err = handleConnection(conn)
			if err != nil {
				logger.ErrorContext(ctx, err.Error())
				os.Exit(1)
			}
		}
	}()
	wg.Wait()
}

func handleConnection(conn quic.Connection) error {
	msg, err := conn.ReceiveMessage(conn.Context())
	if err != nil {
		return err
	}

	if string(msg) == "PING" {
		err := conn.SendMessage([]byte("PONG"))
		if err != nil {
			return err
		}
	}
	return nil
}
