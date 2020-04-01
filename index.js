const { initSession, getSession } = require('./lib/session');
const { initWebsocket } = require('./lib/websocket');

(async function() {
  await initSession();
  initWebsocket();
})();
