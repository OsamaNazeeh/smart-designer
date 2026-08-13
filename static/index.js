const API_URL = "http://localhost:8080";


// ============================================================
// ELEMENTS
// ============================================================

const authPage = document.getElementById("auth-page");
const appPage = document.getElementById("app-page");

const authForm = document.getElementById("auth-form");
const authTitle = document.getElementById("auth-title");
const authButton = document.getElementById("auth-button");
const toggleAuth = document.getElementById("toggle-auth");
const authMessage = document.getElementById("auth-message");

const uploadForm = document.getElementById("upload-form");
const uploadButton = document.getElementById("upload-button");
const uploadMessage = document.getElementById("upload-message");
const preview = document.getElementById("preview");

const transformForm = document.getElementById("transform-form");
const transformButton = document.getElementById("transform-button");
const transformMessage = document.getElementById("transform-message");
const transformedPreview =
    document.getElementById("transformed-preview");

const logoutButton =
    document.getElementById("logout-button");


// ============================================================
// AUTH STATE
// ============================================================

let isRegistering = false;


function getToken() {
    return localStorage.getItem("access_token");
}


// ============================================================
// PAGE VISIBILITY
// ============================================================

function showAuthPage() {
    authPage.classList.remove("hidden");
    appPage.classList.add("hidden");
}


function showAppPage() {
    authPage.classList.add("hidden");
    appPage.classList.remove("hidden");
}


// ============================================================
// LOGOUT / EXPIRED TOKEN
// ============================================================

function logout() {
    localStorage.removeItem("access_token");

    uploadForm.reset();
    transformForm.reset();

    preview.removeAttribute("src");
    transformedPreview.removeAttribute("src");

    preview.classList.add("hidden");
    transformedPreview.classList.add("hidden");

    uploadMessage.textContent = "";
    transformMessage.textContent = "";

    showAuthPage();
}


// ============================================================
// AUTHENTICATED FETCH
// ============================================================

async function authenticatedFetch(url, options = {}) {

    const token = getToken();

    if (!token) {
        showAuthPage();
        throw new Error("Not authenticated");
    }

    const headers = new Headers(options.headers || {});

    headers.set(
        "Authorization",
        `Bearer ${token}`
    );

    const response = await fetch(url, {
        ...options,
        headers
    });

    if (response.status === 401) {
        logout();

        throw new Error(
            "Your session has expired. Please login again."
        );
    }

    return response;
}


// ============================================================
// INITIAL AUTH CHECK
// ============================================================

if (getToken()) {
    showAppPage();
} else {
    showAuthPage();
}


// ============================================================
// LOGIN / REGISTER TOGGLE
// ============================================================

toggleAuth.addEventListener("click", () => {

    isRegistering = !isRegistering;

    authTitle.textContent =
        isRegistering ? "Register" : "Login";

    authButton.textContent =
        isRegistering ? "Register" : "Login";

    toggleAuth.textContent =
        isRegistering
            ? "Already have an account? Login"
            : "Need an account? Register";

    authMessage.textContent = "";
});


// ============================================================
// LOGIN / REGISTER
// ============================================================

authForm.addEventListener("submit", async (event) => {

    event.preventDefault();

    const username =
        document.getElementById("username").value;

    const password =
        document.getElementById("password").value;

    authButton.disabled = true;
    authMessage.textContent = "";

    const endpoint =
        isRegistering ? "/register" : "/login";

    try {

        const response = await fetch(
            API_URL + endpoint,
            {
                method: "POST",

                headers: {
                    "Content-Type": "application/json"
                },

                body: JSON.stringify({
                    username,
                    password
                })
            }
        );

        const data = await response.json();

        if (!response.ok) {
            throw new Error(
                data.message || "Authentication failed"
            );
        }

        const token = data.token;

        if (!token) {
            throw new Error(
                "Server did not return a token"
            );
        }

        localStorage.setItem(
            "access_token",
            token
        );

        showAppPage();

    } catch (error) {

        authMessage.className = "error";
        authMessage.textContent = error.message;

    } finally {

        authButton.disabled = false;
    }
});


// ============================================================
// GET IMAGE OBJECT
// ============================================================

async function getImage(imageId) {

    const response = await authenticatedFetch(
        `${API_URL}/images/${encodeURIComponent(imageId)}`,
        {
            method: "GET"
        }
    );

    if (!response.ok) {
        throw new Error(
            `Failed to get image (${response.status})`
        );
    }

    /*
     * The Go server returns the actual image bytes.
     *
     * We turn those bytes into a browser URL using a Blob.
     */
    const blob = await response.blob();

    const imageURL = URL.createObjectURL(blob);

    return imageURL;
}


// ============================================================
// UPLOAD IMAGE
// ============================================================

uploadForm.addEventListener("submit", async (event) => {

    event.preventDefault();

    const thumbnailFile =
        document.getElementById("thumbnail").files[0];

    if (!thumbnailFile) {
        return;
    }

    const formData = new FormData();

    formData.append(
        "image",
        thumbnailFile
    );

    uploadButton.disabled = true;
    uploadMessage.className = "";
    uploadMessage.textContent = "Uploading...";

    try {

        const response = await authenticatedFetch(
            API_URL + "/images",
            {
                method: "POST",
                body: formData
            }
        );

        const data = await response.json();

        if (response.status == 401) {
            showAuthPage()

        }

        if (!response.ok) {
            throw new Error(
                data.message || "Upload failed"
            );
        }

        uploadMessage.className = "success";
        uploadMessage.textContent =
            "Image uploaded successfully!";

        /*
         * The server should return the image ID.
         *
         * We don't use an S3 URL here.
         *
         * Instead:
         *
         * Browser
         *     ↓
         * GET /images/{id}
         *     ↓
         * Go server
         *     ↓
         * S3
         *     ↓
         * Go server
         *     ↓
         * Browser
         */
        const imageId = data.id;

        if (!imageId) {
            throw new Error(
                "Server did not return an image ID"
            );
        }

        document.getElementById("image-id").value =
            imageId;

        const imageURL = await getImage(imageId);

        preview.src = imageURL;
        preview.classList.remove("hidden");

    } catch (error) {

        uploadMessage.className = "error";
        uploadMessage.textContent = error.message;

    } finally {

        uploadButton.disabled = false;
    }
});


// ============================================================
// TRANSFORM IMAGE
// ============================================================

transformForm.addEventListener("submit", async (event) => {

    event.preventDefault();

    const imageId =
        document.getElementById("image-id").value.trim();

    if (!imageId) {
        return;
    }


    // --------------------------------------------------------
    // Build transformations
    // --------------------------------------------------------

    const transformations = {};


    // Resize

    const resizeWidth =
        document.getElementById("resize-width").value;

    const resizeHeight =
        document.getElementById("resize-height").value;

    if (resizeWidth && resizeHeight) {

        transformations.resize = {
            width: Number(resizeWidth),
            height: Number(resizeHeight)
        };
    }


    // Crop

    const cropWidth =
        document.getElementById("crop-width").value;

    const cropHeight =
        document.getElementById("crop-height").value;

    const cropX =
        document.getElementById("crop-x").value;

    const cropY =
        document.getElementById("crop-y").value;

    if (
        cropWidth &&
        cropHeight &&
        cropX &&
        cropY
    ) {

        transformations.crop = {
            width: Number(cropWidth),
            height: Number(cropHeight),
            x: Number(cropX),
            y: Number(cropY)
        };
    }


    // Rotate

    const rotate =
        document.getElementById("rotate").value;

    if (rotate) {
        transformations.rotate = Number(rotate);
    }


    // Format

    const format =
        document.getElementById("format").value;

    if (format) {
        transformations.format = format;
    }


    // Filters

    const grayscale =
        document.getElementById("grayscale").checked;

    const sepia =
        document.getElementById("sepia").checked;

    if (grayscale || sepia) {

        transformations.filters = {
            grayscale,
            sepia
        };
    }


    const body = {
        transformations
    };


    transformButton.disabled = true;

    transformMessage.className = "";
    transformMessage.textContent =
        "Transforming...";


    try {

        const response = await authenticatedFetch(
            `${API_URL}/images/${encodeURIComponent(imageId)}/transform`,
            {
                method: "POST",

                headers: {
                    "Content-Type": "application/json"
                },

                body: JSON.stringify(body)
            }
        );

        if (response.status == 401) {
            showAuthPage()    
        }
        const data = await response.json();


        if (!response.ok) {
            throw new Error(
                data.message || "Transformation failed"
            );
        }


        transformMessage.className = "success";
        transformMessage.textContent =
            "Image transformed successfully!";


       
        const imageURL = await getImage(imageId);

        transformedPreview.src = imageURL;
        transformedPreview.classList.remove("hidden");

        preview.classList.remove("hidden");

    } catch (error) {

        transformMessage.className = "error";
        transformMessage.textContent =
            error.message;

    } finally {

        transformButton.disabled = false;
    }
});


// ============================================================
// LOGOUT
// ============================================================

logoutButton.addEventListener("click", () => {
    logout();
});