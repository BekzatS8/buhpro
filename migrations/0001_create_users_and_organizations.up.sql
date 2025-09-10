BEGIN;

-- extension for uuid and hashing funcs
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- users
CREATE TABLE IF NOT EXISTS users (
                                     id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(320) UNIQUE,
    phone VARCHAR(32) UNIQUE,
    password_hash TEXT, -- bcrypt
    full_name VARCHAR(255),
    role VARCHAR(32) NOT NULL DEFAULT 'executor', -- client|executor|coach|admin
    status VARCHAR(32) NOT NULL DEFAULT 'active', -- active|suspended|deleted
    metadata JSONB DEFAULT '{}'::jsonb, -- reserved for future fields
    two_fa_enabled BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
    );

CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);
CREATE INDEX IF NOT EXISTS idx_users_phone ON users (phone);
CREATE INDEX IF NOT EXISTS idx_users_role_status ON users (role, status);

-- organizations
CREATE TABLE IF NOT EXISTS organizations (
                                             id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    owner_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(32) NOT NULL, -- 'TOO' | 'IP' | 'REP' (representative)
    bin_iin VARCHAR(20), -- BIN or IIN
    name VARCHAR(512) NOT NULL,
    legal_address TEXT,
    contact JSONB DEFAULT '{}'::jsonb, -- {phone, email, contact_person}
    verification_documents JSONB DEFAULT '[]'::jsonb, -- list of file meta
    status VARCHAR(32) NOT NULL DEFAULT 'pending_verification', -- pending_verification|verified|rejected
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
    );

CREATE INDEX IF NOT EXISTS idx_orgs_status ON organizations (status);
CREATE INDEX IF NOT EXISTS idx_orgs_bin_iin ON organizations (bin_iin);

COMMIT;
