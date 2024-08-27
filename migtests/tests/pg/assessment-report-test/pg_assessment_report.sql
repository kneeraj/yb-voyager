-- Unsupported Datatypes
CREATE TABLE parent_table (
    id SERIAL PRIMARY KEY,
    common_column1 TEXT,
    common_column2 INTEGER
);

CREATE TABLE child_table (
    specific_column1 DATE
) INHERITS (parent_table);

CREATE TABLE Mixed_Data_Types_Table1 (
    id SERIAL PRIMARY KEY,
    point_data POINT,
    snapshot_data TXID_SNAPSHOT,
    lseg_data LSEG,
    box_data BOX
);

CREATE TABLE Mixed_Data_Types_Table2 (
    id SERIAL PRIMARY KEY,
    lsn_data PG_LSN,
    lseg_data LSEG,
    path_data PATH
);

-- GIST Index on point_data column
CREATE INDEX idx_point_data ON Mixed_Data_Types_Table1 USING GIST (point_data);

-- GIST Index on box_data column
CREATE INDEX idx_box_data ON Mixed_Data_Types_Table1 USING GIST (box_data);

CREATE TABLE orders2 (
    id SERIAL PRIMARY KEY,
    order_number VARCHAR(50) UNIQUE DEFERRABLE, --unique constraint deferrable test
    status VARCHAR(50) NOT NULL,
    shipped_date DATE
);

CREATE OR REPLACE FUNCTION prevent_update_shipped_without_date()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'UPDATE' AND NEW.status = 'shipped' AND NEW.shipped_date IS NULL THEN
        RAISE EXCEPTION 'Cannot update status to shipped without setting shipped_date';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create a Constraint Trigger
CREATE CONSTRAINT TRIGGER enforce_shipped_date_constraint
AFTER UPDATE ON orders2
FOR EACH ROW
WHEN (NEW.status = 'shipped' AND NEW.shipped_date IS NULL)
EXECUTE FUNCTION prevent_update_shipped_without_date();

-- Stored Generated Column
CREATE TABLE employees2 (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    full_name VARCHAR(101) GENERATED ALWAYS AS (first_name || ' ' || last_name) STORED,
    Department varchar(50)
);

--For the ALTER TABLE ADD PK on partitioned DDL in schema
CREATE TABLE sales_region (id int, amount int, branch text, region text, PRIMARY KEY(id, region)) PARTITION BY LIST (region);
CREATE TABLE Boston PARTITION OF sales_region FOR VALUES IN ('Boston');
CREATE TABLE London PARTITION OF sales_region FOR VALUES IN ('London');
CREATE TABLE Sydney PARTITION OF sales_region FOR VALUES IN ('Sydney');

--For exclusion constraints
CREATE TABLE public.test_exclude_basic (
    id integer, 
    name text,
    address text
);
ALTER TABLE ONLY public.test_exclude_basic
    ADD CONSTRAINT no_same_name_address EXCLUDE USING btree (name WITH =, address WITH =);


CREATE TABLE test_xml_type(id int, data xml);

INSERT INTO test_xml_type values(1,'<person>
<name>ABC</name>
<age>34</age>
</person>');

INSERT INTO test_xml_type values(2,'<person>
<name>XYZ</name>
<age>36</age>
</person>');


CREATE ROLE test_policy;

CREATE POLICY policy_test_report ON test_xml_type TO test_policy USING (true);

CREATE POLICY policy_test_fine ON public.test_exclude_basic FOR ALL TO PUBLIC USING (id % 2 = 1);

CREATE POLICY policy_test_fine_2 ON public.employees2  USING (id NOT IN (12,123,41241));

CREATE VIEW sales_employees as
select id, first_name,
last_name, full_name
from employees2 where Department = 'sales'
WITH CHECK OPTION;