const prev = 0;
const next = 1;

const test = location.hostname == "localhost";

document.title = "Gallery | James Keveren";

Object.assign(document.body.style, {
	background: "#000",
	color: "#fff"
})

let img = document.createElement("img");
document.body.appendChild(img);
let objectURL = "";
img.addEventListener("load", e => {
	console.log("revoke");
	URL.revokeObjectURL(objectURL);
})

async function moveToImage(direction) {
	let response = await fetch("/gallery/1-06c.jpg");
	let content = await response.blob();
	let objectURL = URL.createObjectURL(content);
	console.log(objectURL);
	img.src = objectURL;
}

if (test) {
	// very basic client side tests
	moveToImage(next);

} else {
	moveToImage(next);
}
