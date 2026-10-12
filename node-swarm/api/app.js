const express = require('express');
const { Pool } = require('pg');

const app = express();
const port = process.env.PORT || 3000;

const pool = new Pool({
  host: process.env.DB_HOST || 'db',
  port: Number(process.env.DB_PORT || 5432),
  database: process.env.DB_NAME || 'xyzdb',
  user: process.env.DB_USER || 'xyzuser',
  password: process.env.DB_PASSWORD || 'xyzpass'
});

app.get('/', (req, res) => {
  res.json({ service: 'XYZ Node.js API', status: 'running' });
});

app.get('/health', (req, res) => {
  res.json({ status: 'healthy' });
});

app.get('/db-check', async (req, res) => {
  try {
    const result = await pool.query('SELECT NOW() AS database_time');
    res.json({
      database: 'connected',
      time: result.rows[0].database_time
    });
  } catch (error) {
    res.status(503).json({ database: 'unavailable' });
  }
});

app.listen(port, '0.0.0.0', () => {
  console.log(`XYZ API listening on port ${port}`);
});
