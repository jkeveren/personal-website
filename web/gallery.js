const prev = 0;
const next = 1;

document.title = "Gallery | James Keveren";

Object.assign(document.body.style, {
	background: "#000",
	color: "#fff",
	margin: 0
})

let img = document.createElement("img");
Object.assign(img.style, {
	position: "fixed",
	maxWidth: "100vw",
	maxHeight: "100vh"
});
document.body.appendChild(img);
let objectURL = "";
img.addEventListener("load", e => {
	URL.revokeObjectURL(objectURL);
});

let response = await fetch("/galleryImage/92708177_p0.jpg");
let content = await response.blob();
objectURL = URL.createObjectURL(content);
console.log(objectURL);
img.src = objectURL;
