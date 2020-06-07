const functions = require('firebase-functions');
const createContentWriter = require('./createContentWriter.js');

const targetHostname = 'james.keve.ren';
const targetProtocol = 'https';
const targetURL = `${targetProtocol}://${targetHostname}/`;

const writeContent = createContentWriter({targetHostname, targetURL});

exports.index = functions.https.onRequest(async (request, response) => {
	try {
		response.setHeader('content-type', 'text/html');
		response.setHeader('cache-control', 'no-store');
		await writeContent({response, targetURL});
	} catch (error) {
		console.error(error);
		response.statusCode = 500;
		response.write('server error');
	} finally {
		response.send();
	}
});
