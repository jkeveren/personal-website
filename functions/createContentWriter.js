const crypto = require('crypto');
const delay = require('./delay.js');

module.exports = ({targetHostname, targetURL}) => {
	const title = 'James Keveren';
	const description = 'Software Developer';
	const email = 'james@keve.ren';

	const createPrecontent = fillerLength =>
	`<html style="font-family:monospace">
	<title>${title}</title>
	<meta name=viewport content=width=device-width,user-scalable=no />
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
	<script src="https://www.gstatic.com/firebasejs/7.6.1/firebase-app.js"></script>
	<script src="https://www.gstatic.com/firebasejs/7.6.1/firebase-analytics.js"></script>
	<script>
		firebase.initializeApp({
			apiKey: "AIzaSyD1kvnPN1m8ree6R-HSKzOlLffhMsEhmHI",
			authDomain: "website-ba4ce.firebaseapp.com",
			databaseURL: "https://website-ba4ce.firebaseio.com",
			projectId: "website-ba4ce",
			storageBucket: "website-ba4ce.appspot.com",
			messagingSenderId: "56181516216",
			appId: "1:56181516216:web:8d173d560ddb0a0ac3d14d",
			measurementId: "G-ZB4SMXQCZ0"
		});
		firebase.analytics();
	</script>
	<!--
		This is just a filler comment to consume a few bytes so browsers start rendering content as it arrives.

		Here's how many bytes it takes some browsers to start rendering content on next byte sent (when sent with "content-type" header of "text/html"):
		- Google Chrome 78.0.3904.70:       2
		- Mozilla Firefox 69.0.3:           1023
		- Microsoft Edge 44.18362.387.0:    511
		- Internet Explorer 11.418.18362.0: 4095

		If the page loaded instantly anyway it's because this is currently deployed to firebase functions which doesn't handle streaming data well so its been disabled.

		Anyway heres a few unnecessarily cryptographically random and HTML-invalid bytes:
	${fillerLength ? crypto.randomBytes(fillerLength) : ''}
	-->
	`;
	const preContentTargetLength = 4095;
	const precontent = createPrecontent(preContentTargetLength - createPrecontent().length);

	const delayMultiplier = 1;
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
				await delay(30 * delayMultiplier);
			}
		}
	};

	const createLinkWriter = (response, writeSlowText) => async (url, text) => {
		response.write(`<a href=${url} target=_blank >`);
		await writeSlowText(text);
		response.write('</a>');
	};

	const pause = () => delay(500 * delayMultiplier);
	const shortPause = () => delay(250 * delayMultiplier);

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
// 			['YouTube',     'https://www.youtube.com/channel/UCsUkIjz__XMM4fYbKGMGmpQ'], // this link is incorrect
			['GitHub',      'https://github.com/jkeveren'],
// 			['Twitter',     'https://twitter.com/JamesKeveren'],
			['Instagram',   'https://instagram.com/jameskeveren'],
			['Thingiverse', 'https://www.thingiverse.com/jkeveren/designs'],
		]) {
			await writeSlowText('- ');
			response.write(`<a target=_blank href=${url}>`);
			await writeSlowText(text);
			response.write('</a><br>');
		}
	};
}
