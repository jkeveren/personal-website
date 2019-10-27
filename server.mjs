import http from 'http';
import createContentWriter from './createContentWriter.mjs';

(async () => {

	const targetHostname = 'james.keve.ren';
	const targetProtocol = 'http';
	const targetURL = `${targetProtocol}://${targetHostname}/`;

	const config = {
		port: 10000
	};

	try {
		Object.assign(config, (await import('./config.mjs')).default);
	} catch (error) {
		if (error.code !== 'ERR_MODULE_NOT_FOUND') {
			throw error;
		}
		console.warn('No config found. Using defaults.');
	}

	if (process.env.PORT) {
		config.port = process.env.PORT;
	}

	if (isNaN(config.port)) {
		throw new Error('port is not a number');
	}

	const writeContent = createContentWriter({targetHostname, targetURL});

	http.createServer(async (req, res) => {
		try {
			if (req.url !== '/' || (req.headers.host !== targetHostname && !process.argv.includes('dev'))) {
				res.statusCode = 307;
				res.setHeader('location', `${targetProtocol}://${targetHostname}/`);
				return res.end('redirecting');
			}
			res.setHeader('content-type', 'text/html');
			await writeContent({res, targetURL});
			res.statusCode = 200;
			return res.end();
		} catch (error) {
			console.error(error);
			res.statusCode = 500;
			res.end('server error');
		}
	}).listen(config.port, () => {
		console.log(`HTTP on port ${config.port}`);
	});

})();
