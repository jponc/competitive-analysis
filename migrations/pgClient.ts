import { Client } from 'pg'

export const getClient = (): Client => {
  return new Client({
    connectionString: process.env.DB_CONN_URL
  })
}
