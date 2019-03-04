/*
 * Copyright (C) 2017 The "MysteriumNetwork/node" Authors.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package tequilapi

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
)

// APIServer interface represents control methods for underlying http api server
type APIServer interface {
	Wait() error
	StartServing() error
	Stop()
	Address() (string, error)
}

type apiServer struct {
	errorChannel  chan error
	handler       http.Handler
	listenAddress string
	listener      net.Listener
}

// NewServer creates http api server for given address port and http handler
func NewServer(address string, port int, handler http.Handler, corsPolicy CorsPolicy) APIServer {
	server := apiServer{
		make(chan error, 1),
		DisableCaching(ApplyCors(handler, corsPolicy)),
		fmt.Sprintf("%s:%d", address, port),
		nil}
	return &server
}

// Stop method stops underlying http server
func (server *apiServer) Stop() {
	if server.listener == nil {
		return
	}
	server.listener.Close()
}

// Wait method waits for http server to finish handling requests (i.e. when Stop() was called)
func (server *apiServer) Wait() error {
	return <-server.errorChannel
}

// Port method returns bind port for given http server (useful when random port is used)
func (server *apiServer) Address() (string, error) {
	if server.listener == nil {
		return "", errors.New("not bound")
	}
	return extractBoundAddress(server.listener)
}

// Port method returns bind port for given http server (useful when random port is used)
// StartServing starts http request serving
func (server *apiServer) StartServing() error {
	var err error
	server.listener, err = net.Listen("tcp", server.listenAddress)
	if err != nil {
		return err
	}
	go server.serve(server.handler)
	return nil
}

func (server *apiServer) serve(handler http.Handler) {
	server.errorChannel <- http.Serve(server.listener, handler)
}

func extractBoundAddress(listener net.Listener) (string, error) {
	addr := listener.Addr()
	parts := strings.Split(addr.String(), ":")
	if len(parts) < 2 {
		return "", errors.New("Unable to locate address: " + addr.String())
	}
	return addr.String(), nil
}
