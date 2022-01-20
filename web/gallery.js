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

// get array of images
let response = await fetch("/galleryImages");
let imagesString = await response.text();
let images = imagesString.split("\n");
let imageIndex = 0;

let abortController = new AbortController()

// get image from server and push to img element
async function displayImage(i, pushState) {
	let image = images[i];
	// If there's nothing to do; do nothing.
	if (image === undefined) {
		return
	}

	// abort any previous fetch to avoid race condition
	abortController.abort();

	imageIndex = i;

	// SPA navigation
	if (pushState) {
		history.pushState({}, document.title, image);
	}

	// fetch image
	abortController = new AbortController()
	let response = await fetch("/galleryImage/" + image, {
		signal: abortController.signal
	});

	// push image to img element
	let content = await response.blob();
	objectURL = URL.createObjectURL(content);
	img.src = objectURL;
}

// SPA navigation
addEventListener("popstate", e => {
	let image = getImageName(location.pathname)
	displayImage(image, false)
});

// arrow keys
addEventListener("keydown", e => {
	if (e.key === "ArrowRight") {
		e.preventDefault()
		displayImage(imageIndex + 1, true);
	} else if (e.key === "ArrowLeft") {
		e.preventDefault()
		displayImage(imageIndex - 1, true);
	}
});

// scroll wheel
addEventListener("wheel", e => {
	e.preventDefault()
	if (e.deltaY > 0) {
		displayImage(imageIndex + 1, true);
	} else {
		displayImage(imageIndex - 1, true);
	}
}, {passive: false}); // passive: true to allow e.preventdefault()


// get image name from path
let parts = location.pathname.split("/");
let URLImage = parts[parts.length - 1];

// find index of image from url
for (let i = 0; i < images.length; i++) {
	let image = images[i];
	if (image == URLImage) {
		imageIndex = i;
	}
}

// display appropriate image
await displayImage(imageIndex, false);
