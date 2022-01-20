// I haven't written unit tests for this yet. This is mostly because I don't
// have time right now given my unfamiliarity of unit testing client side JS.

// steal my code
console.info("https://github.com/jkeveren/personal-website")

document.title = "Gallery | James Keveren";

// put margin in bin
Object.assign(document.body.style, {
	margin: 0
})

function displayError(err) {
	// TODO: write to element
	console.error(err);
}

// create img element and style it
let img = document.createElement("img");
Object.assign(img.style, {
	position: "fixed",
	width: "100vw",
	height: "100vh"
});
document.body.appendChild(img);
img.addEventListener("error", e => {
	displayError(e)
})

// handle object URLs
let objectURL = "";
img.addEventListener("load", e => {
	URL.revokeObjectURL(objectURL);
});

// get array of images
let response = await fetch("/galleryImages");
let imagesString = await response.text();
let images = imagesString.split("\n");

// get imageIndex of image from URL
function getURLImageIndex() {
	// get image name from path
	let parts = location.pathname.split("/");
	let URLImage = parts[parts.length - 1];
	
	let URLImageIndex = null;
	// find index of image from url
	for (let i = 0; i < images.length; i++) {
		let image = images[i];
		if (image == URLImage) {
			URLImageIndex = i;
		}
	}
	if (URLImageIndex === null) {
		return [null, new Error("image not found")]
	}
	return [URLImageIndex, null]
}

// get image from server and push to img element
let imageIndex = 0;
async function displayImage(i, pushState) {
	let image = images[i];
	// If there's nothing to do; do nothing.
	if (image === undefined) {
		return
	}

	imageIndex = i;

	// SPA navigation
	if (pushState) {
		history.pushState({}, document.title, image);
	}

	// Firefox refuses to abort correctly regardless of request technique. I've
	// tried fetch with AbortController, XMHttpRequest with it's abort method and
	// just setting the src value of an img element and none of them abort
	// correctly on firefox. The aborting does remove the image load race
	// condition caused by skipping multiple images but it still continues to
	// load a good portion of the image in the background which consumes
	// bandwidth.
	//
	// Setting src to "" first clears the image so the next image is displayed
	// progressively as it loads. Good for slow connections.
	img.src = "";
	img.src = "/galleryImage/" + image;
}

// SPA navigation
addEventListener("popstate", async e => {
	let [i, err] = getURLImageIndex()
	if (err != null) {
		displayError(err)
	}
	await displayImage(i, false)
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
}, {passive: false}); // passive: false to allow e.preventdefault()

// display image from URL
let [i = 0, err] = getURLImageIndex()
if (err != null) {
	displayError(err)
}
await displayImage(i, false);
