const log4js = require('log4js');
const config = require('../config.json');
const pkg = require('../package.json');

const logger = log4js.getLogger(pkg.name);
logger.level = config.log.level || 'info';

module.exports = {
  logger,
};
