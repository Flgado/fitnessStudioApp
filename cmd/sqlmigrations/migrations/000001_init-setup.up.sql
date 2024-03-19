CREATE TABLE classes (
    id SERIAL PRIMARY KEY,
    class_name VARCHAR(50) NOT NULL,
    class_date TIMESTAMP WITH TIME ZONE NOT NULL,
    class_capacity INT NOT NULL,
    num_registrations INT NOT NULL DEFAULT 0,
    create_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_update_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    user_name VARCHAR(50) NOT NULL,
    create_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_update_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE booking (
    user_id INT REFERENCES users(id),
    class_id INT REFERENCES classes(id),
    reserved_date TIMESTAMP WITH TIME ZONE NOT NULL,
    PRIMARY KEY (user_id, class_id)
);

-- Triggers --
CREATE OR REPLACE FUNCTION update_classes_last_update_date()
RETURNS TRIGGER AS $$
BEGIN
    NEW.last_update_date = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER classes_last_update_trigger
BEFORE UPDATE ON classes
FOR EACH ROW
EXECUTE FUNCTION update_classes_last_update_date();


CREATE OR REPLACE FUNCTION update_users_last_update_date()
RETURNS TRIGGER AS $$
BEGIN
    NEW.last_update_date = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER users_last_update_trigger
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_users_last_update_date();



INSERT INTO classes (class_name, class_date, class_capacity)
VALUES ('CROSSFIT', '2024-04-15', 10);

INSERT INTO users (user_name)
VALUES ('Joao Folgado'), ('Sergio Folgado');


