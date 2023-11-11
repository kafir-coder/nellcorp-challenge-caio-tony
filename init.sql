DROP DATABASE IF EXISTS nell_challenge_test;
CREATE DATABASE nell_challenge_test;

CREATE TABLE IF NOT EXISTS public.accounts (
  id VARCHAR(36) PRIMARY KEY,
  owner_id VARCHAR(36) NOT NULL,
  balance MONEY NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  deletedAt TIMESTAMP NULL
);


CREATE TABLE IF NOT EXISTS public.transaction_ (
  id VARCHAR(36) PRIMARY KEY,
  from_account VARCHAR(36) NOT NULL,
  to_account VARCHAR(36) NOT NULL,
  tr_status TEXT NOT NULL,
  operation VARCHAR(10) NOT NULL,
  amount MONEY NOT NULL,
  multiBeneficiaryId VARCHAR(36) NOT NULL,
  is_Refund BOOLEAN NOT NULL, 
  refunded_transaction_id VARCHAR(36) NOT NULL,
  createdAt TIMESTAMP NOT NULL
);

