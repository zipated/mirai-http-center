const Ajv = require('ajv');
const axios = require('axios');
const { logger } = require('./logger');
const config = require('../config.json');

const ajv = new Ajv();

function post(postURL, data) {
  axios({
    method: 'post',
    url: postURL,
    data,
  })
    .then(function(res) {
      logger.info(`Post data to '${postURL}' succeed.`);
      logger.debug(res);
    })
    .catch(function(error) {
      logger.error(`Post data to '${postURL}' failed. ${error.message}`);
      logger.debug(error);
    });
}

function handleMessage(dataStr) {
  try {
    const data = JSON.parse(dataStr);
    config.schemas.forEach(function(schema) {
      if (ajv.validate(schema.schema, data) && schema.postURL) {
        logger.info(`Post data to '${schema.postURL}'.`);
        logger.debug(schema.schema);
        logger.debug(data);
        post(schema.postURL, data);
      }
    });
  } catch (error) {
    logger.warn(`Handle message erred. ${error.message}`);
    logger.debug(error);
  }
}

module.exports = {
  handleMessage,
};
