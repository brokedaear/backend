// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package service

//
// import (
// 	"crypto/rand"
// 	"encoding/base32"
// 	"errors"
// 	"time"
// )
//
// func NewSessionToken() string {
// 	bytes := make([]byte, 15)
// 	rand.Read(bytes)
// 	sessionId := base32.StdEncoding.EncodeToString(bytes)
//
// 	return sessionId
// }
//
// const sessionExpiresIn = 30 * 24 * time.Hour
//
// func validateSession(sessionId string) (*Session, error) {
// 	session, ok := getSessionFromStorage(sessionId)
// 	if !ok {
// 		return nil, errors.New("invalid session id")
// 	}
// 	if time.Now().After(session.ExpiresAt) {
// 		return nil, errors.New("expired session")
// 	}
// 	if time.Now().After(session.expiresAt.Sub(sessionExpiresIn / 2)) {
// 		session.ExpiresAt = time.Now().Add(sessionExpiresIn)
// 	}
// 	return session, nil
// }
