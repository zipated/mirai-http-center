const log4js = require('log4js');
const config = require('../config.json');

const logger = log4js.getLogger();
logger.level = config.log.level || 'info';

module.exports = {
  logger,
};
