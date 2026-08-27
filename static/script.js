const fileListEl = document.querySelector("#fileList");

const btn = document.querySelector("#next");
const template = document.querySelector("#props");

const savedTheme = localStorage.getItem("theme");
const themeForm = document.getElementById("theme");


function getprops() {
    return {
        startLimit: parseInt(template.getAttribute("data-limit")) || 0,
        limit: parseInt(template.getAttribute("data-limit")) || 0,
        counter: parseInt(template.getAttribute("data-counter")) || 0,
        offset: parseInt(template.getAttribute("data-offset")) || 0,
        keywords: template.getAttribute("data-keywords") || "",
        baseButtonURL: template.getAttribute("data-base-url") || "",
        order: null,
        orderAscending: null,
        replaceUrl: null,
    }
}

if (template) {

    let props = getprops();

    if (props.counter <= props.limit && btn) {
        btn.disabled = true; // Disable the button
    }


    props.offset += props.limit;
    props.baseButtonURL = `/append?keywords=${props.keywords}&limit=${props.limit}&offset=${props.offset}`;
    props.replaceUrl = `/search?keywords=${props.keywords}&limit=${props.limit}&offset=0`;

    changeBtnURL(props)
    updateList(props)
    updateTemplate(props)

    const sortParams = new Set(['name', 'size', 'modtime']);

    document.body.addEventListener('click', function (event) {

        const el = event.target.closest("[data-name]")
        const params = [el?.getAttribute("data-name"), el?.getAttribute("data-order")]


        const duplacetes = FindDuplcates(GetIds())
        if (duplacetes.length > 0)  {
            console.log(FindDuplcates(GetIds()))
        }


        if (event.target.id === 'next') {
            props.offset += props.limit;
            if (!props.order) {
                props.baseButtonURL = `/append?keywords=${props.keywords}&limit=${props.limit}&offset=${props.offset}`;
                props.replaceUrl = `/search?keywords=${props.keywords}&limit=${props.offset}&offset=0`;
            } else {
                props.baseButtonURL = `/append?keywords=${props.keywords}&limit=${props.limit}&offset=${props.offset}&order=${props.order}&ascending=${props.orderAscending ? 'true' : 'false'}`;
                props.replaceUrl = `/search?keywords=${props.keywords}&limit=${props.offset}&offset=0&order=${props.order}&ascending=${props.orderAscending ? 'true' : 'false'}`;
            }
            changeBtnURL(props)
            updateList(props)
            updateTemplate(props)
        } else if (sortParams.has(params[0])) {
            props.order = params[0]
            props.orderAscending = params[1]
            props.baseButtonURL = `/append?keywords=${props.keywords}&limit=${props.limit}&offset=${props.offset}&order=${params[0]}&ascending=${params[1] === 'ASC' ? 'true' : 'false'}`;
            props.replaceUrl = `/search?keywords=${props.keywords}&limit=${props.offset}&offset=0&order=${params[0]}&ascending=${params[1] === 'ASC' ? 'true' : 'false'}`;
            updateTemplate(props)
            changeBtnURL(props)
        }
        
        if (props.offset >= props.counter) {
            btn.disabled = true;
        } 
    });
}

let listContainer = document.getElementById("list");
if (listContainer) {
    document.getElementById("sort-name").addEventListener("click", (event) => handleSortClick(event))
    document.getElementById("sort-size").addEventListener("click", (event) => handleSortClick(event))
    document.getElementById("sort-date").addEventListener("click", (event) => handleSortClick(event))
}


function updateTemplate(props) {
    template.setAttribute("hx-get", props.baseButtonURL);
    template.setAttribute("data-offset", props.offset);
    template.setAttribute("data-order", props.order);
    htmx.process(template)
}

function changeBtnURL(props) {
    btn.setAttribute("hx-get", props.baseButtonURL);
    btn.setAttribute("hx-replace-url", props.replaceUrl);
    btn.setAttribute("data-offset", props.offset);
    htmx.process(btn)
}

function updateList(props) {
    const { keywords, offset, ...rest } = props
    const baseURL = `/append?keywords=${keywords}&limit=${offset}&offset=0`;
    const sortButtons = [
        {
            element: document.querySelector(".ascending.name"),
            ascending: true
        },
        {
            element: document.querySelector(".descending.name"),
        },
        {
            element: document.querySelector(".ascending.size"),
            ascending: true
        },
        {
            element: document.querySelector(".descending.size"),
        },
        {
            element: document.querySelector(".ascending.modtime"),
            ascending: true
        },
        {
            element: document.querySelector(".descending.modtime"),
        },

    ]

    const _ = sortButtons.map(item => {
        const name = item.element.getAttribute("data-name")
        const ascdesc = item.element.getAttribute("data-order") === "ASC"
        item.element.setAttribute(
            'hx-get', baseURL + `&order=${name}&ascending=${ascdesc ? 'true' : 'false'}`
        )
        item.element.innerHTML = ascdesc ? "🠉" : "🠋"
        htmx.process(item.element)
    })
}






function handleSortClick(event) {

    listContainer = document.getElementById("list")
    listItems = Array.from(listContainer.querySelectorAll("li.item-list"))

    const target = event.target;
    const sortType = target.getAttribute("data-sort-type");
    const ascending = target.getAttribute("data-ascending") === "true";

    sortList(sortType, ascending);

    // Toggle the ascending/descending flag
    target.setAttribute("data-ascending", !ascending);
}


function FindDuplcates(arr) {
    const results = [];
    // pets.includes("cat") // true

    const seen = new Set();
    for (const value of arr) {
        if (seen.has(value)) {
            results.push(value);
        } else {
            seen.add(value);
        }
    }

    return results;
}

function GetIds() {
    const ids = [];
    const collection = document.querySelectorAll(".item-Id");
    collection.forEach(
        item => ids.push(item.textContent.replaceAll('\n', '').trim())
    )

    return ids;
}

// Utils


if (savedTheme && themeForm) {
    document.documentElement.setAttribute("data-theme", savedTheme);
    // Set the radio button as checked
    const radio = themeForm.querySelector(`input[name="theme"][value="${savedTheme}"]`);
    if (radio) radio.checked = true;
}

if (themeForm) {
    themeForm.oninput = function (e) {
        const value = e.target.value;
        document.documentElement.setAttribute("data-theme", value);
        localStorage.setItem("theme", value);
    };
}

const previewSectionEl = document.getElementById("preview-section");

function hidePreview() {
    previewSectionEl.style.display = "none";
}

function deletePreview() {
    const p = document.getElementById("preview");
    while (p.firstChild) {
        p.removeChild(p.firstChild);
}
    previewSectionEl.style.display = "none";
}

function showPreview() {
    previewSectionEl.style.display = "block";
}

function showImage(e) {
    console.log("Showing image:", e);

    const image = document.createElement("img");
    image.src = e;
    image.style.maxWidth = "100%";
    image.style.height = "auto";

    const previewImg = document.getElementById("preview-image");
    previewImg.appendChild(image);
    previewImg.style.display = "block";


}

if (previewSectionEl) {
        
    document.addEventListener("keydown", function (event) {
            if (event.key === "Escape" || event.key === "Esc") {
                    if (typeof previewSectionEl !== "undefined" && previewSectionEl.style.display !== "none") {
                            deletePreview();
                        }
                    }
                });
                
}

            

// theme.oninput = e => {
//     document.firstElementChild.setAttribute('data-theme', e.target.value)
// }

const themeToggleBtn = document.getElementById("theme-toggle");
const sun = "☀";
const moon = "☽"

document.addEventListener("DOMContentLoaded", () => {
    const savedTheme = localStorage.getItem("theme") || "light"; // Default to light theme if not set
    setTheme(savedTheme); // Set the initial theme

    if (themeToggleBtn) {
        if (savedTheme === "dark" && themeToggleBtn !== null) {
            themeToggleBtn.textContent = sun;
        } else {
            themeToggleBtn.textContent = moon;
        }
    }
})

if (themeToggleBtn) {
    

    themeToggleBtn.addEventListener("click", (e) => {

        const root = document.firstElementChild;
        let currentTheme = root.getAttribute("data-theme");

        if (currentTheme === "dark") {
            root.setAttribute('data-theme', 'light')
            themeToggleBtn.textContent = moon
        } else {
            root.setAttribute('data-theme', 'dark')
            themeToggleBtn.textContent = sun
        }

        currentTheme = currentTheme === "light" ? "dark" : "light";
        localStorage.setItem("theme", currentTheme);

    })
}







function setTheme(theme) {
    document.documentElement.setAttribute('data-theme', theme);
}


function sortList(parameter, ascending = true) {
    listItems.sort((a, b) => {
        const aValue = a.dataset[parameter]; // Get the data attribute value
        const bValue = b.dataset[parameter];

        switch (parameter) {
            case "name":
                return ascending ? aValue.localeCompare(bValue) : bValue.localeCompare(aValue);
            case "size":
                return ascending ? parseInt(aValue) - parseInt(bValue) : parseInt(bValue) - parseInt(aValue);
            case "moddate":
                let a = aValue.replace(/\s+[A-Z]+$/, "")
                let b = bValue.replace(/\s+[A-Z]+$/, "")
                return ascending ? new Date(a) - new Date(b) : new Date(b) - new Date(a);
            default:
                return 0;
        }
    });

    // Use a DocumentFragment to append all sorted items at once
    const fragment = document.createDocumentFragment();
    listItems.forEach(item => fragment.appendChild(item)); // Add sorted items to the fragment

    // Clear the current list and append the fragment
    listContainer.innerHTML = ""; // Clear the current list
    listContainer.appendChild(fragment); // Append all items at once
}








