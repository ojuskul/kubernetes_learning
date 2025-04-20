const express = require('express');
const app = express();
const port = 3000;

app.get('/', (req, res) => res.send('Success!'));
app.get('/fail', (req, res) => res.status(500).send('Failure!'));

app.listen(port, () => console.log(`App running on port ${port}`));