// I haven't written unit tests for this yet. This is mostly because I don't
// have time right now given my unfamiliarity of unit testing client side JS.

console.info("https://github.com/jkeveren/personal-website")

document.title = "Gallery | James Keveren";

Object.assign(document.body.style, {
	background: "#000",
	color: "#fff",
	margin: 0
})

let img = document.createElement("img");
Object.assign(img.style, {
	position: "fixed",
	width: "100vw",
	height: "100vh"
});
document.body.appendChild(img);
let objectURL = "";
img.addEventListener("load", e => {
	URL.revokeObjectURL(objectURL);
});

function getImageName(path) {
	let parts = path.split("/");
	return parts[parts.length - 1];
}

let imageName = getImageName("/gallery/testImage.jpg")
let want = "testImage.jpg";
console.assert(imageName === want, `Want: ${want}, Got: ${imageName}`)

let nextImage = null;
let prevImage = null;

async function displayImage(image, pushState) {
	if (image === null) {
		return
	}
	if (pushState) {
		history.pushState({}, document.title, image);
	}
	let response = await fetch("/galleryImage/" + image);
	nextImage = response.headers.get("Next");
	prevImage = response.headers.get("Previous");
	let content = await response.blob();
	objectURL = URL.createObjectURL(content);
	img.src = objectURL;
}

addEventListener("popstate", e => {
	let image = getImageName(location.pathname)
	displayImage(image, false)
});

addEventListener("keydown", e => {
	if (e.key === "ArrowRight") {
		displayImage(nextImage, true);
	} else if (e.key === "ArrowLeft") {
		displayImage(prevImage, true);
	}
});

let image = getImageName(location.pathname);
await displayImage(image, false);
