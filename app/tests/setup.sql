-- SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
--
-- SPDX-License-Identifier: Apache-2.0

CREATE TABLE app_user (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE
);

CREATE TABLE user_session (
    id TEXT NOT NULL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES app_user(id),
    expires_at TIMESTAMPTZ NOT NULL,
);

CREATE TABLE token (
	token STRING NOT NULL UNIQUE,
	expires_at INTEGER NOT NULL,
	user_id INTEGER NOT NULL,

	FOREIGN KEY (user_id) REFERENCES user(id)
)
