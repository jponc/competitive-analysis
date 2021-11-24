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
};

export default migrations;
