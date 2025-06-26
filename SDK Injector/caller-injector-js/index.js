const express = require('express');
const Injector = require("./injector.js")
const bodyParser = require('body-parser');
const axios = require('axios');
const { createLogger, format, transports } = require("winston");

const app = express();
app.use(bodyParser.json());

//let injectorURL = process.env.INJECTOR_URL || 'http://injector.default.svc.cluster.local';
const serviceid = process.env.SERVICE_ID || 'hello';

const logger = createLogger({
    level: 'info',
    format: format.combine(
        format.timestamp(),
        format.printf(({ timestamp, level, message }) => {
        return `${new Date(timestamp).toISOString()},CALLER,${level.toUpperCase()},${message}`;
    })),
    transports: [
        new transports.Console()
    ]
});

const injector = new Injector();

app.post('/', async (req, res) => {
    
    const p = req.body;

    //console.log("Received request body:", p);

    if (!p || typeof p.message !== 'string') {
        return res.status(400).send('Invalid JSON');
    }

    const tsMillis = parseInt(p.message, 10);
    if (isNaN(tsMillis)) {
        return res.status(400).send('Invalid timestamp');
    }

    const start = Date.now();
    

    try {
        /**const resp = await axios.get(`${injectorURL}/services/hello`);
        if (resp.status !== 200) {
            return res.status(404).send('Service not found');
        }**/
        const svc = await injector.getServiceById(serviceid);
        if (!svc) {
            return res.status(404).send('Service not found');
        }
        const end = Date.now();
       

        logger.info(`Service retrieved in ${(end - start)} ms`);

        //const svc = resp.data;

        const targetStart = Date.now();
        

        const targetResp = await axios.get(svc.ServiceAddress);
        const targetEnd = Date.now();
        

        logger.info(`Service invoked in ${(targetEnd - targetStart)} ms`);

        const finish = Date.now();
        const total_latency = finish - tsMillis;
        logger.info(`Total latency is ${total_latency} ms`);

        res.send(`Response from service:\n${targetResp.data}`);
    } catch (error) {
        logger.error(error.message);
        res.status(500).send('Failed to reach injector or call target');
    }
});

const PORT = 8080;
app.listen(PORT, () => {
    logger.info(`Function invoker running on :${PORT}`);
});

