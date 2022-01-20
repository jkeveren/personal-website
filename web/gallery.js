// I haven't written unit tests for this yet. This is mostly because I don't
// have time right now given my unfamiliarity of unit testing client side JS.

// steal my code
console.info("https://github.com/jkeveren/personal-website")

document.title = "Gallery | James Keveren";

// put margin in bin
Object.assign(document.body.style, {
	margin: 0
})

// create img element and style it
let img = document.createElement("img");
Object.assign(img.style, {
	position: "fixed",
	width: "100vw",
	height: "100vh"
});
document.body.appendChild(img);

// handle object URLs
let objectURL = "";
img.addEventListener("load", e => {
	URL.revokeObjectURL(objectURL);
});

let nextImage = null;
let prevImage = null;
let abortController = new AbortController()

// get image from server and push to img element
async function displayImage(image, pushState) {
	// If there's nothing to do; do nothing.
	if (image === null) {
		return
	}

	// SPA navigation
	if (pushState) {
		history.pushState({}, document.title, image);
	}

	// abort any previous fetch to avoid race condition
	abortController.abort()
	abortController = new AbortController()
	let response = await fetch("/galleryImage/" + image, {
		signal: abortController.signal
	});

	// get next and prev image
	nextImage = response.headers.get("Next");
	prevImage = response.headers.get("Previous");

	// push image to img element
	let content = await response.blob();
	objectURL = URL.createObjectURL(content);
	img.src = objectURL;
}

// get image name from path
function getImageName(path) {
	let parts = path.split("/");
	return parts[parts.length - 1];
}

// SPA navigation
addEventListener("popstate", e => {
	let image = getImageName(location.pathname)
	displayImage(image, false)
});

// arrow keys
addEventListener("keydown", e => {
	e.preventDefault()
	if (e.key === "ArrowRight") {
		displayImage(nextImage, true);
	} else if (e.key === "ArrowLeft") {
		displayImage(prevImage, true);
	}
});

// scroll wheel
addEventListener("wheel", e => {
	e.preventDefault()
	if (e.deltaY > 0) {
		displayImage(nextImage, true);
	} else {
		displayImage(prevImage, true);
	}
}, {passive: false}); // passive: true to allow e.preventdefault()

// display first image
let image = getImageName(location.pathname);
await displayImage(image, false);
