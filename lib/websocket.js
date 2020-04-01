const axios = require('axios');
const WebSocket = require('ws');
const { logger } = require('./log');
const { getSession } = require('./session');
const config = require('../config.json');

let ws;

async function enableWebsocket() {
  try {
    const enableWebsocketRes = await axios({
      method: 'post',
      url: '/config',
      baseURL: config.mirai.apiBaseURL,
      data: {
        sessionKey: getSession(),
        enableWebsocket: true,
      },
    });
    const { data: enableWebsocketData } = enableWebsocketRes;
    if (enableWebsocketData.code === 0) {
      logger.info('Enable websocket succeed.');
    } else {
      logger.warn('Enable websocket failed.');
    }
    logger.debug(enableWebsocketData);
  } catch (error) {
    logger.fatal(`Enable websocket failed. ${error.message}`);
    logger.debug(error);
    process.exit(-1);
  }
}

function onMessageHandle(data) {
  if (config.websocket.postURL) {
    axios({
      method: 'post',
      url: config.websocket.postURL,
      data: JSON.parse(data),
    })
      .then(function(res) {
        logger.info('Post data succeed.');
        logger.debug(res);
      })
      .catch(function(error) {
        logger.error(`Post data failed. ${error.message}`);
        logger.debug(error);
      });
  }
}

async function initWebsocket() {
  await enableWebsocket();
  ws = new WebSocket(`${config.websocket.baseURL}/all?sessionKey=${getSession()}`);
  ws.on('open', function() {
    logger.info('Websocket opened.');
  });
  ws.on('message', function(data) {
    logger.info('Receive message from websocket.');
    logger.debug(data);
    onMessageHandle(data);
  });
  ws.on('error', function(error) {
    logger.error('Websocket erred.');
    logger.debug(error);
  });
  ws.on('close', function(code, reason) {
    logger.fatal(`Websocket closed with code ${code} and reason '${reason}'.`);
    process.exit(-1);
  });
}

module.exports = {
  initWebsocket,
};
