ALTER TABLE refreshtokens
ADD CONSTRAINT unique_user_id UNIQUE (user_id);
