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
  v06_add_zenserp_batch_processed: async (client: Client) => {
    await client.query(`
      ALTER TABLE query_job ADD COLUMN zenserp_batch_processed BOOLEAN DEFAULT FALSE;
    `);
  },
  v07_add_error_processing_bool_to_query_item: async (client: Client) => {
    await client.query(`
      ALTER TABLE query_item ADD COLUMN error_processing BOOLEAN DEFAULT FALSE;
    `);
  },
  v08_remove_not_null_query_item_body: async (client: Client) => {
    await client.query(`
      ALTER TABLE query_item ALTER COLUMN body DROP NOT NULL;
    `);
  },
  v09_create_link: async (client: Client) => {
    await client.query(`
      CREATE TABLE link
        (
           id             UUID DEFAULT uuid_generate_v4(),
           query_item_id  UUID NOT NULL,
           text           UUID NOT NULL,
           url            TEXT NOT NULL,
           PRIMARY KEY(id),
           CONSTRAINT fk_query_item FOREIGN KEY(query_item_id) REFERENCES query_item(id)
        );
    `);
  },
  v10_update_text_to_text: async (client: Client) => {
    await client.query(`
      ALTER TABLE link ALTER COLUMN text TYPE TEXT;
    `);
  },
  v11_query_jobs_created_at_idx: async (client: Client) => {
    await client.query(`
      CREATE INDEX query_job_created_at_desc_idx ON query_job (created_at DESC);
    `);
  },
  v12_add_query_item_idx: async (client: Client) => {
    await client.query(`
      CREATE INDEX query_item_qj_id_url_idx ON query_item (query_job_id, url);
    `);
  },
  v13_add_query_item_id_idx: async (client: Client) => {
    await client.query(`
      CREATE INDEX link_query_item_id_idx ON link (query_item_id);
    `);
  },
  v14_add_delete_cascade_link: async (client: Client) => {
    await client.query(`
      ALTER TABLE link
      DROP CONSTRAINT fk_query_item,
      ADD CONSTRAINT fk_query_item FOREIGN KEY (query_item_id) REFERENCES query_item(id) ON DELETE CASCADE;
    `);
  },
  v15_add_delete_cascade_query_item_location: async (client: Client) => {
    await client.query(`
      ALTER TABLE query_item
      DROP CONSTRAINT fk_query_location,
      ADD CONSTRAINT fk_query_location FOREIGN KEY (query_location_id) REFERENCES query_location(id) ON DELETE CASCADE;
    `);
  },
  v16_add_delete_cascade_query_item_job: async (client: Client) => {
    await client.query(`
      ALTER TABLE query_item
      DROP CONSTRAINT fk_query_job,
      ADD CONSTRAINT fk_query_job FOREIGN KEY (query_job_id) REFERENCES query_job(id) ON DELETE CASCADE;
    `);
  },
  v17_add_delete_cascade_query_location_job: async (client: Client) => {
    await client.query(`
      ALTER TABLE query_location
      DROP CONSTRAINT fk_query_job,
      ADD CONSTRAINT fk_query_job FOREIGN KEY (query_job_id) REFERENCES query_job(id) ON DELETE CASCADE;
    `);
  },
};

export default migrations;
