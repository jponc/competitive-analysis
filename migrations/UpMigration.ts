import { Client } from "pg";
import { getClient } from "./pgClient";
import migrations from "./list";

export const handler = async () => {
  console.log("Running migrations...")

  const dbClient = getClient();
  await dbClient.connect()

  const currentNames = await getCurrentMigratedNames(dbClient)
  const allNames = Object.keys(migrations);

  for (let i = 0; i < allNames.length; i++) {
    const name = allNames[i];
    if (!currentNames.includes(name)) {
      console.log(`Migrating: ${name}`);
      try {
        await migrations[name](dbClient)
        await addToMigration(dbClient, name)

      } catch (e) {
        console.log(`Failed: ${e.message}`);
        return "Failed";
      }

    }

  }
}

const getCurrentMigratedNames = async (client: Client): Promise<string[]> => {
  const migrationTableQuery = `
    CREATE TABLE IF NOT EXISTS migration (
      created_at TIMESTAMP DEFAULT NOW(),
      name       VARCHAR(100) NOT NULL
    )
  `;
  await client.query(migrationTableQuery);
  const res = await client.query<MigrationRow>('SELECT * FROM migration ORDER BY created_at');
  const migrations = res.rows;

  return migrations.map(m => m.name)
}

type MigrationRow = {
  created_at: Date;
  name: string;
}

const addToMigration = async (client: Client, name: string): Promise<void> => {
  const query = 'INSERT INTO migration(name, created_at) VALUES($1, NOW())'
  const values = [name]

  await client.query(query, values);
}
