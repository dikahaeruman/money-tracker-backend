CREATE TABLE IF NOT EXISTS users(
                                    id serial PRIMARY KEY,
                                    username VARCHAR(50) UNIQUE NOT NULL,
                                    password VARCHAR(300) NOT NULL,
                                    email VARCHAR(300) UNIQUE NOT NULL,
                                    google_uuid UUID,
                                    apple_uuid UUID,
                                    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                                    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                                    deleted_at TIMESTAMP WITH TIME ZONE
);
