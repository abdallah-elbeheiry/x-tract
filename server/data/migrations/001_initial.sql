-- Base User Table
CREATE TABLE IF NOT EXISTS users (
                                     id UUID PRIMARY KEY,
                                     name TEXT NOT NULL,
                                     email TEXT UNIQUE NOT NULL,
                                     hash BYTEA NOT NULL,
                                     salt BYTEA NOT NULL,
                                     number TEXT,
                                     role TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS customer (
                                        user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
                                        stats TEXT
);

CREATE TABLE IF NOT EXISTS salesman (
                                        user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
                                        stats TEXT
);

CREATE TABLE IF NOT EXISTS groups (
                                      id UUID PRIMARY KEY,
                                      name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS guest_employee (
                                              user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
                                              group_id UUID REFERENCES groups(id)
);

CREATE TABLE IF NOT EXISTS operations (
                                          id UUID PRIMARY KEY,
                                          salesman_id UUID REFERENCES salesman(user_id),
                                          customer_id UUID REFERENCES customer(user_id),
                                          location TEXT,
                                          weight NUMERIC,
                                          occurred_at TIMESTAMP WITH TIME ZONE
);