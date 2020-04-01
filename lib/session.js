const axios = require('axios');
const { logger } = require('./log');
const config = require('../config.json');

let session;

async function auth() {
  try {
    const authRes = await axios({
      method: 'post',
      url: '/auth',
      baseURL: config.mirai.apiBaseURL,
      data: {
        authKey: config.mirai.authKey,
      },
    });
    const { data: authData } = authRes;
    if (authData.code === 0) {
      session = authData.session;
      logger.info('Auth succeed.');
      logger.debug(authData);
    } else {
      logger.fatal(`Auth failed with code ${authData.code}.`);
      logger.debug(authData);
      process.exit(-1);
    }
  } catch (error) {
    logger.fatal(`Auth failed. ${error.message}`);
    logger.debug(error);
    process.exit(-1);
  }
}

async function verify() {
  try {
    const verifyRes = await axios({
      method: 'post',
      url: '/verify',
      baseURL: config.mirai.apiBaseURL,
      data: {
        sessionKey: session,
        qq: config.mirai.qq,
      },
    });
    const { data: verifyData } = verifyRes;
    if (verifyData.code === 0) {
      logger.info('Verify session succeed.');
      logger.debug(verifyData);
    } else {
      logger.fatal(`Verify session failed with code ${verifyData.code}.`);
      logger.debug(verifyData);
      process.exit(-1);
    }
  } catch (error) {
    logger.fatal(`Verify session failed. ${error.message}`);
    logger.debug(error);
    process.exit(-1);
  }
}

async function initSession() {
  await auth();
  await verify();
}

function getSession() {
  return session;
}

module.exports = {
  initSession,
  getSession,
};
