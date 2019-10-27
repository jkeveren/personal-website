const crypto = require('crypto');
const delay = require('./delay.js');

module.exports = ({targetHostname, targetURL}) => {
	const title = 'James Keveren';
	const description = 'Software Developer';
	const email = 'james@keve.ren';

	const createPrecontent = fillerLength =>
	`<!DOCTYPE html>
	<html style="font-family:roboto mono;font-size:0.9em">
	<title>${title}</title>
	<meta name=viewport content=width=device-width,user-scalable=no />
	<link href="https://fonts.googleapis.com/css?family=Roboto+Mono&display=swap" rel="stylesheet">
	<meta name="title" content="${title}">
	<meta name="description" content="${description}">
	<meta property="og:type" content="website">
	<meta property="og:url" content="${targetURL}">
	<meta property="og:title" content="${title}">
	<meta property="og:description" content="${description}">
	<meta property="twitter:url" content="${targetURL}">
	<meta property="twitter:title" content="${title}">
	<meta property="twitter:description" content="${description}">
	<link rel="icon" type="image/png" href="data:image/png">
	<!-- Global site tag (gtag.js) - Google Analytics -->
	<script async src="https://www.googletagmanager.com/gtag/js?id=UA-107575308-1"></script>
	<script>
	  (window.dataLayer = window.dataLayer || []).push('js', new Date(), 'config', 'UA-107575308-1');
	</script>
	<!--
		This is just a filler comment to consume a few bytes so browsers start rendering content as it arrives.

		Here's how many bytes it takes some browsers to start rendering content immediately (when sent with "content-type" header of "text/html"):
		- Google Chrome 78.0.3904.70:       3
		- Mozilla Firefox 69.0.3:           1024
		- Microsoft Edge 44.18362.387.0:    512
		- Internet Explorer 11.418.18362.0: 4096

		You can get away with using one less padding byte becasue the first byte of content causes the browser to "render" the invisible padding bytes and the content byte itself.
		Not that a one byte difference is important. Just interesting to me because im weird.

		Anyway heres a few unnecessarily cryptographically random and HTML-invalid bytes (generated at server start for efficiency of course):
	${fillerLength ? crypto.randomBytes(fillerLength) : ''}
	-->\n`;
	const preContentTargetLength = 4095;
	const precontent = createPrecontent(preContentTargetLength - createPrecontent().length);

	const delayMultiplier = 50;

	const writeText = async (response, text) => {
		for (const character of text) {
			response.write(character);
		}
	};

	// replaces newlines with br elements and writes text slowly and br elements instantly. also removes tab characters to allow for code indentation.
	const createSlowTextWriter = response => async text => {
		const tokens = text.replace(/\t/g, '').split(/\n/g);
		for (const lineIndex in tokens) {
			const line = tokens[lineIndex];
			if (lineIndex !== '0') {
				response.write('<br>');
			}
			for (const character of line) {
				response.write(character);
				await delay(30);
			}
		}
	};

	const createLinkWriter = (response, writeSlowText) => async (url, text) => {
		response.write(`<a href=${url}>`);
		await writeSlowText(text);
		response.write('</a>');
	};

	const pause = () => delay(500);
	const shortPause = () => delay(250);

	return async ({response}) => {
		const writeSlowText = createSlowTextWriter(response);
		const writeLink = createLinkWriter(response, writeSlowText);

		// Name + contact
		response.write(precontent);
		await writeSlowText(
			`James Keveren
			`
		);
		await writeLink(targetURL, targetHostname);
		response.write('<br>');
		await writeLink(`mailto:${email}`, email);
		response.write('<br>');

		// Subheading
		await pause();
		await writeSlowText(
			`
			I write software and make things.
			`
		);

// 		// Links
		await pause();
		await writeSlowText(
			`
			Links:
			`
		);
		await shortPause();
		for (const [text, url] of [
// 			['YouTube',     'https://www.youtube.com/channel/UCsUkIjz__XMM4fYbKGMGmpQ'],
			['GitHub',      'https://github.com/jkeveren'],
// 			['Twitter',     'https://twitter.com/JamesKeveren'],
			['Instagram',   'https://instagram.com/jameskeveren'],
			['Thingiverse', 'https://thingiverse.com/jkeveren'],
		]) {
			await writeSlowText('- ');
			response.write(`<a target=_blank href=${url}>`);
			await writeSlowText(text);
			response.write('</a><br>');
		}
	
		// Fin.
		await pause();
		await writeSlowText(
			`
			Fin.`
		);
	};
}