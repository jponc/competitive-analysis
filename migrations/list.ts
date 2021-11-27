import { Client } from "pg";

const migrations = {
  v00_add_uuid_extension: async (client: Client) => {
    await client.query(`
      CREATE EXTENSION "uuid-ossp";
    `);
  },
  v01_create_category_info: async (client: Client) => {
    await client.query(`
      CREATE TABLE query_job
        (
           id            UUID DEFAULT uuid_generate_v4(),
           keyword       TEXT NOT NULL,
           created_at    TIMESTAMP NOT NULL,
           completed_at  TIMESTAMP NOT NULL,
           PRIMARY KEY(id)
        );
    `);
  },
  v01_remove_not_null: async (client: Client) => {
    await client.query(`
      ALTER TABLE query_job ALTER COLUMN completed_at DROP NOT NULL;
    `);
  },
  v02_set_created_at_default: async (client: Client) => {
    await client.query(`
      ALTER TABLE query_job ALTER COLUMN created_at SET DEFAULT NOW();
    `);
  },
  v03_create_query_location: async (client: Client) => {
    await client.query(`
      CREATE TABLE query_location
        (
           id            UUID DEFAULT uuid_generate_v4(),
           query_job_id  UUID NOT NULL,
           device        TEXT NOT NULL,
           search_engine TEXT NOT NULL,
           num           INTEGER NOT NULL,
           country       TEXT NOT NULL,
           location      TEXT NOT NULL,
           created_at    TIMESTAMP NOT NULL DEFAULT NOW(),
           PRIMARY KEY(id),
           CONSTRAINT fk_query_job FOREIGN KEY(query_job_id) REFERENCES query_job(id)
        );
    `);
  },
  v04_create_query_item: async (client: Client) => {
    await client.query(`
      CREATE TABLE query_item
        (
           id                 UUID DEFAULT uuid_generate_v4(),
           query_job_id       UUID NOT NULL,
           query_location_id  UUID NOT NULL,
           position           INTEGER NOT NULL,
           url                TEXT NOT NULL,
           title              TEXT NOT NULL,
           body               TEXT NOT NULL,
           processed_at       TIMESTAMP,
           created_at         TIMESTAMP NOT NULL DEFAULT NOW(),
           PRIMARY KEY(id),
           CONSTRAINT fk_query_job FOREIGN KEY(query_job_id) REFERENCES query_job(id),
           CONSTRAINT fk_query_location FOREIGN KEY (query_location_id) REFERENCES query_location (id)
        );
    `);
  },
  v05_add_zenserp_batch_id_to_query_job: async (client: Client) => {
    await client.query(`
      ALTER TABLE query_job ADD COLUMN zenserp_batch_id TEXT;
    `);
  },
  v06_update_num_to_text: async (client: Client) => {
    await client.query(`
      ALTER TABLE query_location ALTER COLUMN num TYPE TEXT;
    `);
  },
};

export default migrations;
