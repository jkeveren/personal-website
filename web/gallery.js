// I haven't written unit tests for this yet. This is mostly because I don't
// have time right now given my unfamiliarity of unit testing client side JS.

// use this for handling errors that have no recourse
function fatalError(err) {
	// TODO: display error to user
	console.error(err)
}

// yoink my code
console.info("https://github.com/jkeveren/personal-website")

document.title = "Gallery | James Keveren";

// clean up body
Object.assign(document.body.style, {
	margin: 0,
	background: "rgb(127,127,127)", // 50% grey is least jarring
	fontFamily: "sans-serif",
	overflow: "hidden", // disable scollbars
})

// create img element
let img = document.createElement("img");
Object.assign(img.style, {
	width: "100vw",
	height: "100vh",
	objectFit: "contain", // This is what scales the image correctly
});
document.body.appendChild(img);
img.addEventListener("error", e => {
	fatalError(e)
})

// get array of images
let response = await fetch("/galleryImages");
let imagesString = await response.text();
let images = imagesString.split("\n");
let imageIndex = 0;

// create on screen controls
// Swiping is better on a touch screen device but I want to stick with the low
// tech theme of the home page. I'm not styling the buttons either. Enjoy those
// default html buttons.
let buttons = [];
for (let i = 0; i < 2; i++) {
	let b = document.createElement("button");
	buttons[i] = b;
	let height = "60px"
	Object.assign(b.style, {
		position: "fixed",
		width: "50px",
		height: height,
		top: `calc(50vh - ${height}/2)`
	});
	if (i === 0) {
		b.style.left = 0;
		b.textContent = "<";
		b.title = "Previous Image";
		b.addEventListener("click", e => {
			displayImage(imageIndex - 1, true);
		});
	} else {
		b.style.right = 0;
		b.textContent = ">";
		b.title = "Next Image";
		b.addEventListener("click", e => {
			displayImage(imageIndex + 1, true);
		});
	}
	document.body.appendChild(b);
}

const imageNotFoundError = new Error("image not found");
const emptyURLImageError = new Error("no image name was found in URL");

// get imageIndex of image from URL
function getURLImageIndex() {
	// get image name from path
	let parts = location.pathname.split("/");
	// use first image if no image is specified
	let URLImage = parts[parts.length - 1];
	if (URLImage === "") {
		return [0, emptyURLImageError];
	}
	let URLImageIndex = null;
	// find index of image from url
	for (let i = 0; i < images.length; i++) {
		let image = images[i];
		if (image == URLImage) {
			URLImageIndex = i;
		}
	}
	if (URLImageIndex === null) {
		return [null, imageNotFoundError]
	}
	return [URLImageIndex, null]
}

// get image from server and push to img element
async function displayImage(i, pushState) {
	let image = images[i];
	// If there's nothing to do; do nothing.
	if (image === undefined) {
		return
	}

	imageIndex = i;

	// manage button state
	buttons[0].disabled = i === 0;
	buttons[1].disabled = i === images.length - 1;

	// SPA navigation
	if (pushState) {
		history.replaceState({}, document.title, image);
	}

	// Firefox refuses to abort correctly regardless of request technique. I've
	// tried fetch with AbortController, XMHttpRequest with it's abort method and
	// just setting the src value of an img element and none of them abort
	// correctly on firefox. The aborting does remove the image load race
	// condition caused by skipping multiple images but it still continues to
	// load a good portion of the image in the background which consumes
	// bandwidth.
	// 
	// Using css background-image: url() does not abort previous requests at all
	// in any browser that I tried.
	//
	// Setting src to "" first clears the image so the next image is displayed
	// progressively as it loads. Good for slow connections.
	img.src = "";
	img.src = "/galleryImage/" + image;
}

// SPA navigation
addEventListener("popstate", async e => {
	displayImageFromURL();
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
async function displayImageFromURL() {
	let [i = 0, err] = getURLImageIndex()
	if (err != null && err != emptyURLImageError) {
		fatalError(err)
		throw "throwing because it's impossible to return here"
	}
	await displayImage(i, err === emptyURLImageError);
}

displayImageFromURL();